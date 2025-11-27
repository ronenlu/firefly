package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"firefly/internal/data"
	"firefly/internal/engine"
	"firefly/internal/words"
)

func main() {
	ctx := context.Background()

	urls := data.LoadEssayURLs()
	bank := words.LoadWordBank()

	result, err := engine.RunSequential(ctx, urls, bank)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
