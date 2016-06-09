package ripper

import "github.com/evalphobia/go-jp-text-ripper/tokenizer"

// Plugin outputs extra column with custom logic
type Plugin struct {
	Title string
	Fn    func(rawText, normalizedText string, tokens *tokenizer.TokenList) string
}
