package data

import (
	"bufio"
	_ "embed"
	"strings"
)

//go:embed endg-urls.txt
var urlsFile string

func LoadEssayURLs() []string {
	scanner := bufio.NewScanner(strings.NewReader(urlsFile))
	var urls []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		urls = append(urls, line)
	}
	return urls
}
