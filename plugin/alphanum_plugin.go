package plugin

import (
	"regexp"
	"strconv"

	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

var (
	reAlphaNum = regexp.MustCompile(`[a-zA-Z0-9]`)
	reAlphabet = regexp.MustCompile(`[a-zA-Z]`)
	reNumber   = regexp.MustCompile(`[0-9]`)
)

// AlphaNumCountPlugin calculates alphabet and number count from normalized text
var AlphaNumCountPlugin = &ripper.Plugin{
	Title: "alphanum_count",
	Fn: func(text *ripper.TextData) string {
		count := len(reAlphaNum.FindAllString(text.GetNormalized(), -1))
		return strconv.Itoa(count)
	},
}

// AlphabetCountPlugin calculates alphabet count from normalized text
var AlphabetCountPlugin = &ripper.Plugin{
	Title: "alphabet_count",
	Fn: func(text *ripper.TextData) string {
		count := len(reAlphabet.FindAllString(text.GetNormalized(), -1))
		return strconv.Itoa(count)
	},
}

// NumberCountPlugin calculates Number count from normalized text
var NumberCountPlugin = &ripper.Plugin{
	Title: "number_count",
	Fn: func(text *ripper.TextData) string {
		count := len(reNumber.FindAllString(text.GetNormalized(), -1))
		return strconv.Itoa(count)
	},
}
