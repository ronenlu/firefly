package words

import (
	"bufio"
	_ "embed"
	"regexp"
	"strings"
	"unicode"
)

//go:embed word_bank.txt
var wordBankFile string

type Bank map[string]struct{}

func LoadWordBank() Bank {
	scanner := bufio.NewScanner(strings.NewReader(wordBankFile))
	bank := make(map[string]struct{})
	for scanner.Scan() {
		w := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if w == "" {
			continue
		}
		bank[w] = struct{}{}
	}
	return bank
}

var nonLetters = regexp.MustCompile(`[^A-Za-z]+`)

// CountValidWords counts occurrences of valid words in the given text based on the provided word bank.
func CountValidWords(text string, bank Bank) map[string]int {
	text = strings.ToLower(text)
	tokens := nonLetters.Split(text, -1)

	counts := make(map[string]int)
	for _, t := range tokens {
		if !isValidWord(t, bank) {
			continue
		}
		counts[t]++
	}
	return counts
}

func isValidWord(w string, bank Bank) bool {
	if len(w) < 3 {
		return false
	}
	for _, r := range w {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	if _, ok := bank[w]; !ok {
		return false
	}
	return true
}
