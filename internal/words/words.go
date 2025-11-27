package words

import (
	"bufio"
	_ "embed"
	"strings"
)

//go:embed word_bank.txt
var wordBankFile string

func LoadWordBank() map[string]struct{} {
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
