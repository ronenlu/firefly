package engine

import (
	"context"
	"firefly/internal/words"
	"io"
	"net/http"
	"runtime"
	"sort"
	"sync"
)

type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type Output struct {
	Top []WordCount `json:"top_words"`
}

func RunConcurrent(ctx context.Context, urls []string, bank words.Bank) (*Output, error) {
	maxWorkers := 64
	numWorkers := runtime.NumCPU() * 4
	if numWorkers > maxWorkers {
		numWorkers = maxWorkers
	}

	// Channels for work distribution and results
	urlChan := make(chan string, len(urls))
	resultChan := make(chan map[string]int, len(urls))
	errChan := make(chan error, len(urls))

	// Global word counts with mutex for safe concurrent access
	global := make(map[string]int)
	var mu sync.Mutex

	// WaitGroup to track worker completion
	var wg sync.WaitGroup

	// HTTP client shared across workers
	client := &http.Client{}

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range urlChan {
				// Check context cancellation
				select {
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				default:
				}

				// Fetch URL
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
				if err != nil {
					errChan <- err
					return
				}

				resp, err := client.Do(req)
				if err != nil {
					errChan <- err
					return
				}

				body, err := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				if err != nil {
					errChan <- err
					return
				}

				// Process word counts
				local := words.CountValidWords(string(body), bank)
				resultChan <- local
			}
		}()
	}

	// Send URLs to workers
	go func() {
		for _, u := range urls {
			urlChan <- u
		}
		close(urlChan)
	}()

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	// Collect results
	var firstErr error
	for local := range resultChan {
		mu.Lock()
		merge(global, local)
		mu.Unlock()
	}

	// Check for errors
	for err := range errChan {
		if firstErr == nil {
			firstErr = err
		}
	}

	if firstErr != nil {
		return nil, firstErr
	}

	// top 10 + pretty JSON
	top := topN(global, 10)
	out := Output{Top: top}
	return &out, nil
}

func merge(global, local map[string]int) {
	for w, c := range local {
		global[w] += c
	}
}

func topN(counts map[string]int, n int) []WordCount {
	var all []WordCount
	for w, c := range counts {
		all = append(all, WordCount{Word: w, Count: c})
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].Count == all[j].Count {
			return all[i].Word < all[j].Word
		}
		return all[i].Count > all[j].Count
	})
	if len(all) > n {
		all = all[:n]
	}
	return all
}
