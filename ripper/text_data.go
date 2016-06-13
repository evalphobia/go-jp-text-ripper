package ripper

import "github.com/evalphobia/go-jp-text-ripper/tokenizer"

type TextData struct {
	raw        string
	normalized string
	words      *tokenizer.TokenList
	nonWords   *tokenizer.TokenList

	Optional string // optional field for plugins
}

func (t *TextData) GetRaw() string {
	return t.raw
}

func (t *TextData) GetNormalized() string {
	return t.normalized
}

func (t *TextData) GetWords() *tokenizer.TokenList {
	return t.words
}

func (t *TextData) GetNonWords() *tokenizer.TokenList {
	return t.nonWords
}
