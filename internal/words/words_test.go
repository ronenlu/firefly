package words

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountValidWords(t *testing.T) {
	// Create a test bank with known words
	testBank := Bank{
		"the":     struct{}{},
		"quick":   struct{}{},
		"brown":   struct{}{},
		"fox":     struct{}{},
		"jumps":   struct{}{},
		"over":    struct{}{},
		"lazy":    struct{}{},
		"dog":     struct{}{},
		"hello":   struct{}{},
		"world":   struct{}{},
		"testing": struct{}{},
		"golang":  struct{}{},
		"code":    struct{}{},
	}

	tests := []struct {
		name string
		text string
		bank Bank
		want map[string]int
	}{
		{
			name: "simple sentence",
			text: "The quick brown fox jumps over the lazy dog",
			bank: testBank,
			want: map[string]int{
				"the":   2,
				"quick": 1,
				"brown": 1,
				"fox":   1,
				"jumps": 1,
				"over":  1,
				"lazy":  1,
				"dog":   1,
			},
		},
		{
			name: "repeated words",
			text: "hello hello world world world",
			bank: testBank,
			want: map[string]int{
				"hello": 2,
				"world": 3,
			},
		},
		{
			name: "mixed case",
			text: "HELLO Hello hello WORLD World world",
			bank: testBank,
			want: map[string]int{
				"hello": 3,
				"world": 3,
			},
		},
		{
			name: "with punctuation",
			text: "Hello, world! Testing... testing; golang.",
			bank: testBank,
			want: map[string]int{
				"hello":   1,
				"world":   1,
				"testing": 2,
				"golang":  1,
			},
		},
		{
			name: "words not in bank",
			text: "invalid unknown missing words",
			bank: testBank,
			want: map[string]int{},
		},
		{
			name: "short words (less than 3 chars)",
			text: "the a an is to be",
			bank: testBank,
			want: map[string]int{
				"the": 1,
			},
		},
		{
			name: "empty text",
			text: "",
			bank: testBank,
			want: map[string]int{},
		},
		{
			name: "only punctuation",
			text: "!@#$%^&*()",
			bank: testBank,
			want: map[string]int{},
		},
		{
			name: "numbers and special characters",
			text: "hello123 world456 code789",
			bank: testBank,
			want: map[string]int{
				"hello": 1,
				"world": 1,
				"code":  1,
			},
		},
		{
			name: "mixed valid and invalid",
			text: "The quick xyz jumps abc over",
			bank: testBank,
			want: map[string]int{
				"the":   1,
				"quick": 1,
				"jumps": 1,
				"over":  1,
			},
		},
		{
			name: "empty bank",
			text: "hello world testing",
			bank: Bank{},
			want: map[string]int{},
		},
		{
			name: "whitespace variations",
			text: "  hello   world  \t testing  \n golang  ",
			bank: testBank,
			want: map[string]int{
				"hello":   1,
				"world":   1,
				"testing": 1,
				"golang":  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountValidWords(tt.text, tt.bank)
			assert.Equal(t, tt.want, got)
		})
	}
}
