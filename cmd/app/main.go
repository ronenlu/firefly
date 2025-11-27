package main

import (
	"context"
	"log"

	"firefly/internal/data"
	"firefly/internal/engine"
	"firefly/internal/words"
)

func main() {
	ctx := context.Background()

	urls := data.LoadEssayURLs()
	bank := words.LoadWordBank()

	if err := engine.RunSequential(ctx, urls, bank); err != nil {
		log.Fatal(err)
	}
}
