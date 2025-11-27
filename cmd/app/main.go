package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"firefly/internal/data"
	"firefly/internal/engine"
	"firefly/internal/words"
)

func main() {
	ctx := context.Background()

	urls := data.LoadEssayURLs()
	bank := words.LoadWordBank()

	start := time.Now()
	result, err := engine.RunConcurrent(ctx, urls, bank)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("time took for processing:", time.Since(start))

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
