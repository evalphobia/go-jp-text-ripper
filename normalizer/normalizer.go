package normalizer

import "strings"

// Normalizer formats raw text to normalized text
type Normalizer struct {
	Normalize func(rawText string) string
}

// Default is default *Normalizer
var Default = &Normalizer{
	Normalize: func(rawText string) string {
		return defaultReplacer.Replace(rawText)
	},
}

var defaultReplacer = strings.NewReplacer(`↵`, " ",
	`"`, " ",
	`　`, " ",
	`\t`, " ",
	"\t", " ",
	`\n`, " ",
	"\n", " ")
