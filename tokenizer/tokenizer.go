package tokenizer

import (
	"github.com/ikawaha/kagome/tokenizer"
)

const (
	PosNoun      = "名詞"
	PosVerb      = "動詞"
	PosAdjective = "形容詞"
)

var defaultWordPosList = []string{
	PosNoun,
	PosVerb,
	PosAdjective,
}

// Tokenizer is struct for tokenize text
type Tokenizer struct {
	t tokenizer.Tokenizer

	minLetterSize   int
	wordPosList     []string
	wordPosMap      map[string]struct{}
	stopWordMap     map[string]struct{}
	useOriginalForm bool
}

// New returns initialized Tokenizer.
func New(c Config) *Tokenizer {
	t := &Tokenizer{
		t:               tokenizer.New(),
		wordPosList:     defaultWordPosList,
		minLetterSize:   1,
		useOriginalForm: c.UseOriginalForm,
	}

	if c.MinLetterSize > 1 {
		t.minLetterSize = c.MinLetterSize
	}

	if len(c.WordPosList) != 0 {
		t.wordPosList = c.WordPosList
	}
	t.wordPosMap = make(map[string]struct{}, len(t.wordPosList))
	for _, p := range t.wordPosList {
		t.wordPosMap[p] = struct{}{}
	}

	t.stopWordMap = make(map[string]struct{}, len(c.StopWordList))
	for _, p := range c.StopWordList {
		t.stopWordMap[p] = struct{}{}
	}

	return t
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

// AddStopWords adds word into stop word list.
func (t *Tokenizer) AddStopWords(list ...string) {
	for _, p := range list {
		t.stopWordMap[p] = struct{}{}
	}
}

// Tokenize separates text into tokens(words) and return the list
func (t *Tokenizer) Tokenize(text string) (*TokenList, *TokenList) {
	tokens := t.t.Tokenize(text)

	words := make([]*Token, 0, len(tokens))
	nonWords := make([]*Token, 0, len(tokens))
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			continue
		}

		nt := newToken(token)
		if t.isValidWord(nt.GetPos(), nt.GetSurface()) {
			words = append(words, nt)
		} else {
			nonWords = append(nonWords, nt)
		}
	}

	wordList := &TokenList{
		List:            words,
		UseOriginalForm: t.useOriginalForm,
	}
	nonList := &TokenList{
		List:            nonWords,
		UseOriginalForm: t.useOriginalForm,
	}
	return wordList, nonList
}

func (t *Tokenizer) isValidWord(pos, surface string) bool {
	if _, ok := t.wordPosMap[pos]; !ok {
		return false
	}
	if len(surface) < t.minLetterSize {
		return false
	}
	if _, ok := t.stopWordMap[surface]; ok {
		return false
	}
	// ignore a word which letters contains only special signs.
	if !ReCJKPatterns.MatchString(surface) {
		return false
	}

	return true
}

// Config for Tokenizer.
type Config struct {
	MinLetterSize   int
	WordPosList     []string
	StopWordList    []string
	UseOriginalForm bool
}
