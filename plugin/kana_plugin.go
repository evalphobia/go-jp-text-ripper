package plugin

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// KanaNumberLikePlugin calculates number-like japanese word count from normalized text
var KanaNumberLikePlugin = &ripper.Plugin{
	Title: "kana_number_count",
	Fn: func(text *ripper.TextData) string {
		count := strings.Count(jpNumberReplacer.Replace(strings.ToLowerSpecial(kanaConv, text.GetNormalized())), jpSymbol)
		return strconv.Itoa(count)
	},
}

// KanaAlphabetLikePlugin calculates alphabet-like japanese word count from normalized text
var KanaAlphabetLikePlugin = &ripper.Plugin{
	Title: "kana_alphabet_count",
	Fn: func(text *ripper.TextData) string {
		count := strings.Count(jpAlphabetReplacer.Replace(strings.ToLowerSpecial(kanaConv, text.GetNormalized())), jpSymbol)
		return strconv.Itoa(count)
	},
}

var kanaConv = unicode.SpecialCase{
	unicode.CaseRange{
		Lo: 0x30a1, // ア
		Hi: 0x30f3, // ン
		Delta: [unicode.MaxCase]rune{
			0,
			0x3041 - 0x30a1,
			0,
		},
	},
}

const jpSymbol = "\x01"

var jpNumberReplacer = strings.NewReplacer(
	`ぜろ`, jpSymbol,
	`いち`, jpSymbol,
	// "に", jpSymbol,
	"さん", jpSymbol,
	"よん", jpSymbol,
	// "ご", jpSymbol,
	"ろく", jpSymbol,
	"しち", jpSymbol,
	"はち", jpSymbol,
	"きゅう", jpSymbol,
	"きゅー", jpSymbol,
	"じゅう", jpSymbol,
	"じゅー", jpSymbol,
	"わん", jpSymbol,
	"つー", jpSymbol,
	"すりー", jpSymbol,
	"ふぉー", jpSymbol,
	"ふぁいぶ", jpSymbol,
	"しっくす", jpSymbol,
	"せぶん", jpSymbol,
	"えいと", jpSymbol,
	"ないん", jpSymbol,
	"てん", jpSymbol)

var jpAlphabetReplacer = strings.NewReplacer(
	`えー`, jpSymbol,
	`びー`, jpSymbol, `びい`, jpSymbol,
	`しー`, jpSymbol,
	`でぃー`, jpSymbol,
	`いー`, jpSymbol,
	`えふ`, jpSymbol,
	`じー`, jpSymbol,
	`えっち`, jpSymbol, `えいち`, jpSymbol,
	`あい`, jpSymbol,
	`じぇー`, jpSymbol, `じぇい`, jpSymbol,
	`けー`, jpSymbol, `けい`, jpSymbol,
	`える`, jpSymbol,
	`えむ`, jpSymbol,
	`えす`, jpSymbol,
	`おー`, jpSymbol,
	`ぴー`, jpSymbol,
	`きゅー`, jpSymbol,
	`あーる`, jpSymbol,
	`えす`, jpSymbol,
	`てぃー`, jpSymbol,
	`ゆー`, jpSymbol,
	`ぶい`, jpSymbol,
	`だぶる`, jpSymbol, `だぶりゅ`, jpSymbol,
	`えっくす`, jpSymbol,
	`わい`, jpSymbol,
	`ぜっと`, jpSymbol)
