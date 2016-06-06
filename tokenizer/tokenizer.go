package tokenizer

import (
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

// Tokenizer tokenize text
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

func (t *Tokenizer) SetDictinary(path string) error {
	dic, err := tokenizer.NewDic(path)
	if err != nil {
		return err
	}

	t.t.SetDic(dic)
	return nil
}

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

func (t *Tokenizer) prefilter(text string) string {
	return t.replacer.Replace(text)
}

func newReplacer() *strings.Replacer {
	return strings.NewReplacer("↵", " ",
		"\t", " ",
		"\n", " ")
}
