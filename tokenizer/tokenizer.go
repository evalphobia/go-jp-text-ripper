package tokenizer

import "github.com/ikawaha/kagome/tokenizer"

// Tokenizer is struct for tokenize text
type Tokenizer struct {
	t tokenizer.Tokenizer
}

// New returns initialized Tokenizer
func New() *Tokenizer {
	return &Tokenizer{
		t: tokenizer.New(),
	}
}

// SetDictionary sets new dictionary for tokenize
func (t *Tokenizer) SetDictionary(path string) error {
	dic, err := tokenizer.NewDic(path)
	if err != nil {
		return err
	}

	t.t.SetDic(dic)
	return nil
}

// Tokenize separates text into tokens(words) and return the list
func (t *Tokenizer) Tokenize(text string) *TokenList {
	tokens := t.t.Tokenize(text)

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
