package tokenizer

import (
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

// Tokenizer is struct for tokenize text
type Tokenizer struct {
	t        tokenizer.Tokenizer
	replacer *strings.Replacer
}

// New returns initialized Tokenizer
func New() *Tokenizer {
	return &Tokenizer{
		t:        tokenizer.New(),
		replacer: newReplacer(),
	}
}

// SetDictinary sets new dictinary for tokenize
func (t *Tokenizer) SetDictinary(path string) error {
	dic, err := tokenizer.NewDic(path)
	if err != nil {
		return err
	}

	t.t.SetDic(dic)
	return nil
}

// Tokenize separates text into tokens(words) and return the list
func (t *Tokenizer) Tokenize(text string) *TokenList {
	tokens := t.t.Tokenize(t.prefilter(text))

	result := make([]*Token, 0, len(tokens))
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			continue
		}

		nt := newToken(token)
		if !nt.isWord() {
			continue
		}

		result = append(result, nt)
	}

	list := &TokenList{
		List: result,
	}
	return list
}

// prefilter preprocesses text before tokenize
func (t *Tokenizer) prefilter(text string) string {
	return t.replacer.Replace(text)
}

// newReplacer returns strings.Replacer to remove unwanted symbols
func newReplacer() *strings.Replacer {
	return strings.NewReplacer("↵", " ",
		`\t`, " ",
		`\n`, " ",
		"\t", " ",
		`"`, " ",
		"\n", " ")
}
