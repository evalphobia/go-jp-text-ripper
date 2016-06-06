package tokenizer

import "github.com/ikawaha/kagome/tokenizer"

// TokenList is token
type TokenList struct {
	List []*Token
}

func (list *TokenList) GetWords() []string {
	words := make([]string, len(list.List))
	for i, t := range list.List {
		words[i] = t.Token.Surface
	}
	return words
}

// Token is token
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
