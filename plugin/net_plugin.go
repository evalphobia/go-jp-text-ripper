package plugin

import (
	"regexp"
	"strings"

	"github.com/evalphobia/go-jp-text-ripper/ripper"
	"github.com/evalphobia/go-jp-text-ripper/tokenizer"
)

var (
	domainRegExp   = regexp.MustCompile(`\.[a-zA-Z]{2,}`)
	domainReplacer = strings.NewReplacer("↵", "",
		`\s`, " ",
		`　`, " ")
)

// DomainLikePlugin calculates domain-like word count from tokenized words
var DomainLikePlugin = &ripper.Plugin{
	Title: "domain_like_count",
	Fn: func(_, normalizedText string, _ *tokenizer.TokenList) string {
		return ""
	},
}

func domainLikeFilter(in string) string {
	return domainReplacer.Replace(in)
}
