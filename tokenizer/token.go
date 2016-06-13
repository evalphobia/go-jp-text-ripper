package tokenizer

import "github.com/ikawaha/kagome/tokenizer"

// Token is a tokenized word
type Token struct {
	tokenizer.Token
	pos string
}

func newToken(token tokenizer.Token) *Token {
	t := &Token{
		Token: token,
		pos:   token.Pos(),
	}
	return t
}

func (t *Token) isNoun() bool {
	return t.pos == "名詞"
}

func (t *Token) isVerb() bool {
	return t.pos == "動詞"
}

func (t *Token) isAdjective() bool {
	return t.pos == "形容詞"
}

func (t *Token) isWord() bool {
	return t.isNoun() || t.isVerb() || t.isAdjective()
}

func (t *Token) GetPos() string {
	return t.pos
}

// HasFeature checks token contains the feature or not
func (t *Token) HasFeature(f string) bool {
	for _, val := range t.Token.Features() {
		if val == f {
			return true
		}
	}
	return false
}

// TokenList is token slice list
type TokenList struct {
	List []*Token
}

// GetWords returns word list
func (list *TokenList) GetWords() []string {
	words := make([]string, len(list.List))
	for i, t := range list.List {
		words[i] = t.Token.Surface
	}
	return words
}

// CountFeatures counts matched feature
func (list *TokenList) CountFeatures(f string) int {
	count := 0
	for _, t := range list.List {
		if t.HasFeature(f) {
			count++
		}
	}
	return count
}
