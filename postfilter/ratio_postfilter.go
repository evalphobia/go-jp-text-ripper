package postfilter

import (
	"strconv"

	"github.com/evalphobia/go-jp-text-ripper/plugin"
	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// RatioAlphaNum calculates alphabet and number ratio from raw text
var RatioAlphaNum = &ripper.PostFilter{
	Title: "ratio_alphanum_count",
	Fn: func(data map[string]string) string {
		return getCharacterRatioFromText(data, plugin.AlphaNumCountPlugin.Title)
	},
}

// RatioAlphabet calculates alphabet ratio from raw text
var RatioAlphabet = &ripper.PostFilter{
	Title: "ratio_alphabet_count",
	Fn: func(data map[string]string) string {
		return getCharacterRatioFromText(data, plugin.AlphabetCountPlugin.Title)
	},
}

// RatioNumber calculates number ratio from raw text
var RatioNumber = &ripper.PostFilter{
	Title: "ratio_number_count",
	Fn: func(data map[string]string) string {
		return getCharacterRatioFromText(data, plugin.NumberCountPlugin.Title)
	},
}

// RatioJP calculates japanese character ratio from raw text
var RatioJP = &ripper.PostFilter{
	Title: "ratio_jp_count",
	Fn: func(data map[string]string) string {
		return getCharacterRatioFromText(data, plugin.KanaCountPlugin.Title)
	},
}

func getCharacterRatioFromText(data map[string]string, target string) string {
	targetStr, ok := data[target]
	if !ok {
		return ""
	}

	targetCount, err := strconv.Atoi(targetStr)
	if err != nil {
		return ""
	}

	charStr, ok := data["raw_char_count"]
	if !ok {
		return ""
	}

	charCount, err := strconv.Atoi(charStr)
	if err != nil {
		return ""
	}
	return strconv.FormatFloat(float64(targetCount)/float64(charCount), 'f', 8, 64)
}
