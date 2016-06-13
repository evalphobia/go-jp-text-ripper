package plugin

import (
	"strconv"

	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// CharTypeCountPlugin calculates character type count from normalized text
var CharTypeCountPlugin = &ripper.Plugin{
	Title: "char_type_count",
	Fn: func(text *ripper.TextData) string {
		m := make(map[rune]struct{})
		for _, s := range text.GetNormalized() {
			m[s] = struct{}{}
		}

		return strconv.Itoa(len(m))
	},
}

// MaxCharCountPlugin calculates maximum character type frequency from normalized text
var MaxCharCountPlugin = &ripper.Plugin{
	Title: "max_char_count",
	Fn: func(text *ripper.TextData) string {
		m := make(map[rune]int)
		count := 0
		for _, s := range text.GetNormalized() {
			if _, ok := m[s]; !ok {
				m[s] = 0
			}
			m[s]++
			if count < m[s] {
				count = m[s]
			}
		}
		return strconv.Itoa(count)
	},
}

// MaxWordCountPlugin calculates maximum word frequency from tokenized words
var MaxWordCountPlugin = &ripper.Plugin{
	Title: "max_word_count",
	Fn: func(text *ripper.TextData) string {
		m := make(map[string]int)
		count := 0
		for _, t := range text.GetWords().List {
			s := t.Surface
			if _, ok := m[s]; !ok {
				m[s] = 0
			}
			m[s]++
			if count < m[s] {
				count = m[s]
			}
		}
		return strconv.Itoa(count)
	},
}

// SymbolCountPlugin calculates symbol word count from tokenized words
var SymbolCountPlugin = &ripper.Plugin{
	Title: "symbol_count",
	Fn: func(text *ripper.TextData) string {
		return strconv.Itoa(text.GetNonWords().CountFeatures("記号"))
	},
}
