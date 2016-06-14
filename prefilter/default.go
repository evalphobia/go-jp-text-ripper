package prefilter

import (
	"strings"

	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// DefaultNormalizer is prefilter to remove white spaces
var DefaultNormalizer = &ripper.PreFilter{
	Fn: func(rawText string) string {
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
