package tokenizer

import (
	"regexp"

	"github.com/ikawaha/kagome/tokenizer"
)

// Token is a tokenized word.
type Token struct {
	tokenizer.Token
	pos      string
	features []string

	WordPosList   []string
	MinLetterSize int
}

func newToken(token tokenizer.Token) *Token {
	t := &Token{
		Token:         token,
		pos:           token.Pos(),
		features:      token.Features(),
		MinLetterSize: 1,
	}
	return t
}

// GetPos returns pos text (the first feature).
func (t *Token) GetPos() string {
	return t.pos
}

// GetSurface returns surface text.
func (t *Token) GetSurface() string {
	return t.Token.Surface
}

// GetOriginalForm returns the original form of surface text.
func (t *Token) GetOriginalForm() string {
	if len(t.features) < 7 {
		return t.GetSurface()
	}

	s := t.features[6]
	switch s {
	case "",
		"*":
		return t.GetSurface()
	default:
		return s
	}
}

// HasFeature checks token contains the feature or not.
func (t *Token) HasFeature(f string) bool {
	for _, val := range t.features {
		if val == f {
			return true
		}
	}
	return false
}

// TokenList is token slice list
type TokenList struct {
	List            []*Token
	UseOriginalForm bool
}

// GetWords returns word list
func (list *TokenList) GetWords() []string {
	words := make([]string, len(list.List))

	useOriginal := list.UseOriginalForm
	for i, t := range list.List {
		switch {
		case useOriginal:
			words[i] = t.GetOriginalForm()
		default:
			words[i] = t.GetSurface()
		}
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

// HasFeatures checks if countain matched feature
func (list *TokenList) HasFeatures(f string) bool {
	for _, t := range list.List {
		if t.HasFeature(f) {
			return true
		}
	}
	return false
}

var (
	// from: https://techlife.cookpad.com/entry/2019/02/20/120219
	ReCJKPatterns            = regexp.MustCompile("[" + CJKPatterns + "]+")
	ReNonCJKPatterns         = regexp.MustCompile("[^" + CJKPatterns + "]+")
	ReNonCJKPatternsWithSign = regexp.MustCompile("[^" + CJKPatternsWithSign + "]+")

	CJKPatterns = `\x{3040}-\x{309F}` + // Hiragana
		`\x{30A0}-\x{30FF}` + // Katakana
		`\x{FF65}-\x{FF9F}` + // Half width Katakana
		`\x{FF10}-\x{FF19}` + // Full width digits
		`\x{FF21}-\x{FF3A}` + // Full width Upper case English Alphabets
		`\x{FF41}-\x{FF5A}` + // Full width Lower case English Alphabets
		`\x{0030}-\x{0039}` + // Half width digits
		`\x{0041}-\x{005A}` + // Half width Upper case English Alphabets
		`\x{0061}-\x{007A}` + // Half width Lower case English Alphabets
		`\x{3190}-\x{319F}` + // Kanbun
		`\x{4E00}-\x{9FFF}` // CJK unified ideographs. kanjis

	CJKPatternsWithSign = CJKPatterns +
		`\x{3001}-\x{3002}` + // '、' '。'
		`\x{0021}` + // '!'
		`\x{FF01}` // '！'
)
