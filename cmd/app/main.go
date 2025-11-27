package main

import (
	"fmt"

	"firefly/internal/data"
	"firefly/internal/words"
)

func main() {
	urls := data.LoadEssayURLs()
	fmt.Printf("loaded %d urls\n", len(urls))

	bank := words.LoadWordBank()
	fmt.Printf("loaded %d words in bank\n", len(bank))
}
