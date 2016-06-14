package plugin

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

var (
	reJP       = regexp.MustCompile(`[\p{Han}\p{Hiragana}\p{Katakana}]`)
	reHiragana = regexp.MustCompile(`[\p{Hiragana}]`)
	reKatakana = regexp.MustCompile(`[\p{Katakana}]`)
	reKanji    = regexp.MustCompile(`[\p{Han}]`)
)

// KanaCountPlugin calculates japanese character count from normalized text
var KanaCountPlugin = &ripper.Plugin{
	Title: "kana_count",
	Fn: func(text *ripper.TextData) string {
		count := len(reJP.FindAllString(text.GetNormalized(), -1))
		return strconv.Itoa(count)
	},
}

// HiraganaCountPlugin calculates japanese hiragana character count from normalized text
var HiraganaCountPlugin = &ripper.Plugin{
	Title: "hiragana_count",
	Fn: func(text *ripper.TextData) string {
		count := len(reHiragana.FindAllString(text.GetNormalized(), -1))
		return strconv.Itoa(count)
	},
}

// KatakanaCountPlugin calculates japanese katakana character count from normalized text
var KatakanaCountPlugin = &ripper.Plugin{
	Title: "katakana_count",
	Fn: func(text *ripper.TextData) string {
		count := len(reKatakana.FindAllString(text.GetNormalized(), -1))
		return strconv.Itoa(count)
	},
}

// KanjiCountPlugin calculates japanese kanji character count from normalized text
var KanjiCountPlugin = &ripper.Plugin{
	Title: "kanji_count",
	Fn: func(text *ripper.TextData) string {
		count := len(reKanji.FindAllString(text.GetNormalized(), -1))
		return strconv.Itoa(count)
	},
}

// KanaAlphaNumLikeCountPlugin calculates alphanum-like japanese word count from normalized text
var KanaAlphaNumLikeCountPlugin = &ripper.Plugin{
	Title: "kana_alphanum_count",
	Fn: func(text *ripper.TextData) string {
		t := jpAlphabetReplacer.Replace(strings.ToLowerSpecial(kanaConv, text.GetNormalized()))
		count := strings.Count(jpNumberReplacer.Replace(t), jpSymbol)
		return strconv.Itoa(count)
	},
}

// KanaNumberLikeCountPlugin calculates number-like japanese word count from normalized text
var KanaNumberLikeCountPlugin = &ripper.Plugin{
	Title: "kana_number_count",
	Fn: func(text *ripper.TextData) string {
		count := strings.Count(jpNumberReplacer.Replace(strings.ToLowerSpecial(kanaConv, text.GetNormalized())), jpSymbol)
		return strconv.Itoa(count)
	},
}

// KanaAlphabetLikeCountPlugin calculates alphabet-like japanese word count from normalized text
var KanaAlphabetLikeCountPlugin = &ripper.Plugin{
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
	"にー", jpSymbol,
	// "に", jpSymbol,
	"さん", jpSymbol,
	"よん", jpSymbol,
	"ごー", jpSymbol,
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
	`えっち`, jpSymbol, `えいち`, jpSymbol, `えち`, jpSymbol,
	`あい`, jpSymbol,
	`じぇー`, jpSymbol, `じぇい`, jpSymbol,
	`けー`, jpSymbol, `けい`, jpSymbol,
	`える`, jpSymbol,
	`えむ`, jpSymbol,
	`えぬ`, jpSymbol,
	`おー`, jpSymbol,
	`ぴー`, jpSymbol,
	`きゅー`, jpSymbol,
	`あーる`, jpSymbol,
	`えす`, jpSymbol,
	`てぃー`, jpSymbol,
	`ゆー`, jpSymbol,
	`ぶい`, jpSymbol,
	`だぶる`, jpSymbol, `だぶりゅ`, jpSymbol,
	`えっくす`, jpSymbol, `えくす`, jpSymbol,
	`わい`, jpSymbol,
	`ぜっと`, jpSymbol,
	`ぜっど`, jpSymbol)
