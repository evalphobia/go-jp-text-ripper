go-jp-text-ripper
====

[![Build Status](https://travis-ci.org/evalphobia/go-jp-text-ripper.svg?branch=master)](https://travis-ci.org/evalphobia/go-jp-text-ripper) [![codecov](https://codecov.io/gh/evalphobia/go-jp-text-ripper/branch/master/graph/badge.svg)](https://codecov.io/gh/evalphobia/go-jp-text-ripper)
 [![GoDoc](https://godoc.org/github.com/evalphobia/go-jp-text-ripper?status.svg)](https://godoc.org/github.com/evalphobia/go-jp-text-ripper)

Separate long text into words and put spaces between ths words.


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

- `-input`: input file path (required)
    - supported file format: `csv`, `tsv`
    - auto detect from the file extension
- `-output`: output file path (required)
    - supported file format: `csv`, `tsv`
    - auto detect from the file extension
    - omit this option if you set `-show` flag
- `-column`: target column name (required)
- `-replace`: replace the column for the result words
    - `false`: output result into added new column `op_text`
    - `true`: output result into column `-column`
- `-show`: print separated words on console
- `-debug`: print word details on console
- `-progress`: print current progress on console (default=30)
- `-quote`: columns name to quote output (comma separated)

# License

Apache License, Version 2.0
