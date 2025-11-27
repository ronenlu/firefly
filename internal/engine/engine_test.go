package engine

import (
	"context"
	"firefly/internal/words"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunConcurrent(t *testing.T) {
	// Create a simple word bank for testing
	bank := words.Bank{
		"hello":   struct{}{},
		"world":   struct{}{},
		"test":    struct{}{},
		"golang":  struct{}{},
		"program": struct{}{},
		"code":    struct{}{},
	}

	t.Run("successful processing of multiple URLs", func(t *testing.T) {
		// Create test servers with different content
		server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Hello world! This is a test. Hello hello world."))
		}))
		defer server1.Close()

		server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Golang is great. Test test test. World world."))
		}))
		defer server2.Close()

		server3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Code program code. Hello golang program."))
		}))
		defer server3.Close()

		urls := []string{
			server1.URL,
			server2.URL,
			server3.URL,
		}

		ctx := context.Background()
		result, err := RunConcurrent(ctx, urls, bank)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Top)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Hello world"))
		}))
		defer server.Close()

		urls := make([]string, 10)
		for i := range urls {
			urls[i] = server.URL
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result, err := RunConcurrent(ctx, urls, bank)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("handles HTTP request errors", func(t *testing.T) {
		urls := []string{"http://invalid-url-that-does-not-exist-12345.com"}

		ctx := context.Background()
		result, err := RunConcurrent(ctx, urls, bank)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("handles empty word bank", func(t *testing.T) {
		emptyBank := words.Bank{}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Hello world test"))
		}))
		defer server.Close()

		urls := []string{server.URL}

		ctx := context.Background()
		result, err := RunConcurrent(ctx, urls, emptyBank)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Top)
	})

	t.Run("handles server errors", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		urls := []string{server.URL}

		ctx := context.Background()
		result, err := RunConcurrent(ctx, urls, bank)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestMerge(t *testing.T) {
	t.Run("merges two word count maps", func(t *testing.T) {
		global := map[string]int{
			"hello": 5,
			"world": 3,
		}
		local := map[string]int{
			"hello": 2,
			"test":  1,
		}

		merge(global, local)

		assert.Equal(t, 7, global["hello"])
		assert.Equal(t, 3, global["world"])
		assert.Equal(t, 1, global["test"])
	})

	t.Run("merges with empty global map", func(t *testing.T) {
		global := map[string]int{}
		local := map[string]int{
			"hello": 2,
			"test":  1,
		}

		merge(global, local)

		assert.Equal(t, 2, global["hello"])
		assert.Equal(t, 1, global["test"])
	})
}
