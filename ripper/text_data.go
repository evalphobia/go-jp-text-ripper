package ripper

import "github.com/evalphobia/go-jp-text-ripper/tokenizer"

// TextData is used in plugins for analyzing text data
type TextData struct {
	raw        string
	normalized string
	words      *tokenizer.TokenList
	nonWords   *tokenizer.TokenList

	Optional string // optional field for plugins
}

// GetRaw returns raw text data
func (t *TextData) GetRaw() string {
	return t.raw
}

// GetNormalized returns normalized text data
func (t *TextData) GetNormalized() string {
	return t.normalized
}

// GetWords returns word tokens
func (t *TextData) GetWords() *tokenizer.TokenList {
	return t.words
}

// GetNonWords returns non-word tokens
func (t *TextData) GetNonWords() *tokenizer.TokenList {
	return t.nonWords
}
