package engine

import (
	"context"
	"firefly/internal/words"
	"io"
	"net/http"
	"sort"
)

type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type Output struct {
	Top []WordCount `json:"top_words"`
}

func RunSequential(ctx context.Context, urls []string, bank words.Bank) (*Output, error) {
	client := &http.Client{}
	global := make(map[string]int)

	for _, u := range urls {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, err
		}

		local := words.CountValidWords(string(body), bank)
		merge(global, local)
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
