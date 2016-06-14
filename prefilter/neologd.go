package prefilter

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// Neologd is prefilter to normalize text by neologd recommended format
var Neologd = &ripper.PreFilter{
	Fn: func(rawText string) string {
		return NormalizeNeologd(rawText)
	},
}

// Most of code logic is from https://github.com/ikawaha/x/neologd/neologd.go
const (
	prolongedSoundMark = '\u30FC'
)

var latinSymbols = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0021, 0x0040, 1},
		{0x005B, 0x0060, 1},
		{0x007B, 0x007E, 1},
	},
}

var neologdReplacer = strings.NewReplacer(
	"０", "0", "１", "1", "２", "2", "３", "3", "４", "4",
	"５", "5", "６", "6", "７", "7", "８", "8", "９", "9",

	"Ａ", "A", "Ｂ", "B", "Ｃ", "C", "Ｄ", "D", "Ｅ", "E",
	"Ｆ", "F", "Ｇ", "G", "Ｈ", "H", "Ｉ", "I", "Ｊ", "J",
	"Ｋ", "K", "Ｌ", "L", "Ｍ", "M", "Ｎ", "N", "Ｏ", "O",
	"Ｐ", "P", "Ｑ", "Q", "Ｒ", "R", "Ｓ", "S", "Ｔ", "T",
	"Ｕ", "U", "Ｖ", "V", "Ｗ", "W", "Ｘ", "X", "Ｙ", "Y",
	"Ｚ", "Z",

	"ａ", "a", "ｂ", "b", "ｃ", "c", "ｄ", "d", "ｅ", "e",
	"ｆ", "f", "ｇ", "g", "ｈ", "h", "ｉ", "i", "ｊ", "j",
	"ｋ", "k", "ｌ", "l", "ｍ", "m", "ｎ", "n", "ｏ", "o",
	"ｐ", "p", "ｑ", "q", "ｒ", "r", "ｓ", "s", "ｔ", "t",
	"ｕ", "u", "ｖ", "v", "ｗ", "w", "ｘ", "x", "ｙ", "y",
	"ｚ", "z",

	//small case
	"ｧ", "ァ", "ｨ", "ィ", "ｩ", "ゥ", "ｪ", "ェ", "ｫ", "ォ",
	"ｬ", "ャ", "ｭ", "ュ", "ｮ", "ョ", "ｯ", "ッ",

	"ｱ", "ア", "ｲ", "イ", "ｳ", "ウ", "ｴ", "エ", "ｵ", "オ",
	"ｶﾞ", "ガ", "ｷﾞ", "ギ", "ｸﾞ", "グ", "ｹﾞ", "ゲ", "ｺﾞ", "ゴ",
	"ｶ", "カ", "ｷ", "キ", "ｸ", "ク", "ｹ", "ケ", "ｺ", "コ",
	"ｻﾞ", "ザ", "ｼﾞ", "ジ", "ｽﾞ", "ズ", "ｾﾞ", "ゼ", "ｿﾞ", "ゾ",
	"ｻ", "サ", "ｼ", "シ", "ｽ", "ス", "ｾ", "セ", "ｿ", "ソ",
	"ﾀﾞ", "ダ", "ﾁﾞ", "ヂ", "ﾂﾞ", "ヅ", "ﾃﾞ", "デ", "ﾄﾞ", "ド",
	"ﾀ", "タ", "ﾁ", "チ", "ﾂ", "ツ", "ﾃ", "テ", "ﾄ", "ト",
	"ﾅ", "ナ", "ﾆ", "ニ", "ﾇ", "ヌ", "ﾈ", "ネ", "ﾉ", "ノ",
	"ﾊﾞ", "バ", "ﾋﾞ", "ビ", "ﾌﾞ", "ブ", "ﾍﾞ", "ベ", "ﾎﾞ", "ボ",
	"ﾊﾟ", "パ", "ﾋﾟ", "ピ", "ﾌﾟ", "プ", "ﾍﾟ", "ペ", "ﾎﾟ", "ポ",
	"ﾊ", "ハ", "ﾋ", "ヒ", "ﾌ", "フ", "ﾍ", "ヘ", "ﾎ", "ホ",
	"ﾏ", "マ", "ﾐ", "ミ", "ﾑ", "ム", "ﾒ", "メ", "ﾓ", "モ",
	"ﾔ", "ヤ", "ﾕ", "ユ", "ﾖ", "ヨ",
	"ﾗ", "ラ", "ﾘ", "リ", "ﾙ", "ル", "ﾚ", "レ", "ﾛ", "ロ",
	"ﾜ", "ワ", "ｦ", "ヲ", "ﾝ", "ン",

	// hyphen
	"\u02D7", "-", "\u058A", "-", "\u2010", "-", "\u2011", "-", "\u2012", "-",
	"\u2013", "-", "\u2043", "-", "\u207B", "-", "\u208B", "-", "\u2212", "-",

	// bar
	"\u2014", string(prolongedSoundMark), // エムダッシュ
	"\u2015", string(prolongedSoundMark), // ホリゾンタルバー
	"\u2500", string(prolongedSoundMark), // 横細罫線
	"\u2501", string(prolongedSoundMark), // 横太罫線
	"\uFE63", string(prolongedSoundMark), // SMALL HYPHEN-MINUS
	"\uFF0D", string(prolongedSoundMark), // 全角ハイフンマイナス
	"\uFF70", string(prolongedSoundMark), // 半角長音記号

	// tilde
	"~", "", "\u223C", "", "\u223E", "", "\u301C", "", "\u3030", "", "\uFF5E", "",

	// zen -> han
	"！", "!", "”", `"`, "＃", "#", "＄", "$", "％", "%",
	"＆", "&", `’`, `'`, "（", "(", "）", ")", "＊", "*",
	"＋", "+", "，", ",", "−", "-", "．", ".", "／", "/",
	"：", ":", "；", ";", "＜", "<", "＞", ">", "？", "?",
	"＠", "@", "［", "[", "￥", "\u00A5", "］", "]", "＾", "^",
	"＿", "_", "｀", "`", "｛", "{", "｜", "|", "｝", "}",
	"　", " ",

	// han -> zen
	"｡", "。", "､", "、", "･", "・", "=", "＝", "｢", "「", "｣", "」",

	"\n", " ", `\n`, " ", "\t", " ", `\t`, " ", "\v", " ", `\v`, " ",
)

// NormalizeNeologd normalizes text
func NormalizeNeologd(s string) string {
	s = neologdReplacer.Replace(s)
	s = eliminateSpace(s)
	s = shurinkProlongedSoundMark(s)
	return s
}

func eliminateSpace(s string) string {
	var (
		b    bytes.Buffer
		prev rune
	)
	for p := 0; p < len(s); {
		c, w := utf8.DecodeRuneInString(s[p:])
		p += w
		if !unicode.IsSpace(c) {
			b.WriteRune(c)
			prev = c
			continue
		}
		for p < len(s) {
			c0, w0 := utf8.DecodeRuneInString(s[p:])
			p += w0
			if !unicode.IsSpace(c0) {
				if unicode.In(prev, unicode.Latin, latinSymbols) &&
					unicode.In(c0, unicode.Latin, latinSymbols) {
					b.WriteRune(' ')
				}
				b.WriteRune(c0)
				prev = c0
				break
			}
		}

	}
	return b.String()
}

func shurinkProlongedSoundMark(s string) string {
	var b bytes.Buffer
	for p := 0; p < len(s); {
		c, w := utf8.DecodeRuneInString(s[p:])
		p += w
		b.WriteRune(c)
		if c != prolongedSoundMark {
			continue
		}
		for p < len(s) {
			c0, w0 := utf8.DecodeRuneInString(s[p:])
			p += w0
			if c0 != prolongedSoundMark {
				b.WriteRune(c0)
				break
			}
		}

	}
	return b.String()
}
