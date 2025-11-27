package words

import (
	"bufio"
	_ "embed"
	"strings"
	"unicode"
)

//go:embed word_bank.txt
var wordBankFile string

// Bank represents a set of valid words used for validation.
type Bank map[string]struct{}

// LoadWordBank loads and parses the embedded word bank file,
// returning a Bank with normalized lowercase words.
func LoadWordBank() Bank {
	scanner := bufio.NewScanner(strings.NewReader(wordBankFile))
	bank := make(map[string]struct{})
	for scanner.Scan() {
		w := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if w == "" {
			continue
		}
		if !isValidWord(w) {
			continue
		}
		bank[w] = struct{}{}
	}
	return bank
}

// splitWords splits a string into slices of consecutive Unicode letters.
func splitWords(s string) []string {
	var words []string
	var current []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			current = append(current, unicode.ToLower(r))
		} else {
			if len(current) > 0 {
				words = append(words, string(current))
				current = current[:0]
			}
		}
	}
	if len(current) > 0 {
		words = append(words, string(current))
	}
	return words
}

// CountValidWords counts occurrences of valid words in the given text based on the provided word bank.
func CountValidWords(text string, bank Bank) map[string]int {
	tokens := splitWords(text)
	counts := make(map[string]int)
	for _, t := range tokens {
		if !isValidWord(t) {
			continue
		}
		if _, exists := bank[t]; !exists {
			continue
		}
		counts[t]++
	}
	return counts
}

// // isValidWord returns true if the word is at least 3 characters long and contains only letters.
func isValidWord(w string) bool {
	if len(w) < 3 {
		return false
	}
	for _, r := range w {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
