go-jp-text-ripper
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Co decov Coverage][11]][12] [![Go Report Card][13]][14] [![Code Climate][19]][20] [![BCH compliance][21]][22] [![Downloads][15]][16]

[1]: https://godoc.org/github.com/evalphobia/go-jp-text-ripper?status.svg
[2]: https://godoc.org/github.com/evalphobia/go-jp-text-ripper
[3]: https://img.shields.io/badge/license-Apache%202-blue
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/evalphobia/go-jp-text-ripper.svg
[6]: https://github.com/evalphobia/go-jp-text-ripper/releases/latest
[7]: https://travis-ci.org/evalphobia/go-jp-text-ripper.svg?branch=master
[8]: https://travis-ci.org/evalphobia/go-jp-text-ripper
[9]: https://coveralls.io/repos/evalphobia/go-jp-text-ripper/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/evalphobia/go-jp-text-ripper?branch=master
[11]: https://codecov.io/github/evalphobia/go-jp-text-ripper/coverage.svg?branch=master
[12]: https://codecov.io/github/evalphobia/go-jp-text-ripper?branch=master
[13]: https://goreportcard.com/badge/github.com/evalphobia/go-jp-text-ripper
[14]: https://goreportcard.com/report/github.com/evalphobia/go-jp-text-ripper
[15]: https://img.shields.io/github/downloads/evalphobia/go-jp-text-ripper/total.svg?maxAge=1800
[16]: https://github.com/evalphobia/go-jp-text-ripper/releases
[17]: https://img.shields.io/github/stars/evalphobia/go-jp-text-ripper.svg
[18]: https://github.com/evalphobia/go-jp-text-ripper/stargazers
[19]: https://codeclimate.com/github/evalphobia/go-jp-text-ripper/badges/gpa.svg
[20]: https://codeclimate.com/github/evalphobia/go-jp-text-ripper
[21]: https://bettercodehub.com/edge/badge/evalphobia/go-jp-text-ripper?branch=master
[22]: https://bettercodehub.com/


`go-jp-text-ripper` separates long text of Japanese into words and put spaces between ths words.


# Quick Usage

## command line
```sh
# install
$ go get github.com/evalphobia/go-jp-text-ripper


# execute
$ go-jp-text-ripper -input ./example/input.csv -output ./example/output.tsv -column text -replace -show -debug

read and write lines...
==============================
すももももももものうち
------
すももももももものうち
------
すもも	名詞,一般,*,*,*,*,すもも,スモモ,スモモ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
うち	名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
------
すもも もも もも うち
==============================
こんにちは、ちゃんはむだよ！ よろしくネ☆ミ
------
こんにちは、ちゃんはむだよ！ よろしくネ☆ミ
------
ちゃん	名詞,接尾,人名,*,*,*,ちゃん,チャン,チャン
はむ	動詞,自立,*,*,五段・マ行,基本形,はむ,ハム,ハム
ミ	名詞,一般,*,*,*,*,ミ,ミ,ミ
------
ちゃん はむ ミ
finish process


# show result
$ cat ./example/output.tsv

user_id	status	text	op_word_count	op_raw_char_count
1	0	すもも もも もも うち	4	11
2	1	ちゃん はむ ミ	3	22
```

## Custome Go App

Import `go-jp-text-ripper` and add plugins.
You can add your custome plugins.

```go
package main

import (
	"github.com/evalphobia/go-jp-text-ripper/plugin"
	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// cli entry point
func main() {
	// prefilters to normalize raw text
	ripper.DefaultPreFilters = []*ripper.PreFilter{
			prefilter.Neologd,
	}

	// plugins
	ripper.DefaultPlugins = []*ripper.Plugin{
		plugin.KanaCountPlugin,
		plugin.AlphaNumCountPlugin,
		plugin.CharTypeCountPlugin,
		plugin.MaxCharCountPlugin,
		plugin.MaxWordCountPlugin,
		plugin.SymbolCountPlugin,
		plugin.NounNameCountPlugin,
		plugin.NounHasFullNamePlugin,
		plugin.NounNumberCountPlugin,
		plugin.KanaNumberLikeCountPlugin,
		plugin.KanaAlphabetLikeCountPlugin,
		plugin.NounLocationCountPlugin,
		plugin.NounOrganizationCountPlugin,
		// MyCustomePlugin,
		&ripper.Plugin{
			Title: "proper_noun_count",
			Fn: func(text *ripper.TextData) string {
				return strconv.Itoa(text.GetWords().CountFeatures("固有名詞"))
			},
		},
	}

	// postfilters running after processed all of the plugins
	ripper.DefaultPostFilters = []*ripper.PostFilter{
		postfilter.RatioJP,
		postfilter.RatioAlphaNum,
	}

	ripper.AutoRun()
}
```

then, build and run!

### Options

| option | description | required | default | example |
|:--|:--|:--|:--:|:--:|:--|
| `-column` | Target column name | yes | - |  |
| `-input` | Input file path. Delimiter is auto-detected by the extention. Supported file formats are `csv`, `tsv`  | yes | - | `list.csv` `file.tsv` |
| `-output` | Output file path. Delimiter is auto-detected by the extention. If you set `-show` flag, you can omit this option. | (yes) | - | `output.csv` `result.tsv` |
| `-replace` | Replace the column of the result words. If it's true, the result will be on column `-column`. If it's false, the result will be added on new column `op_text` | - | `false` | `false` `true` |
| `-show` | Print separated words on console. | - | `false` | `false` `true` |
| `-debug` | Print word details on console. | - | `false` | `false` `true` |
| `-progress` | Intervals to print current progress on console | - | 30 | `false` `true` |
| `-quote` | Columns name to quote for the output result. You can set multiple columns with comma. | - |  | `id,name,data` |

# License

Apache License, Version 2.0
