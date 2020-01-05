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

## Install

```sh
# install
$ go get github.com/evalphobia/go-jp-text-ripper

# or clone and build
# $ git clone --depth 1 https://github.com/evalphobia/go-jp-text-ripper.git
# $ cd ./go-jp-text-ripper
# $ make build
```

```sh
$ go-jp-text-ripper -h

Commands:

  help   show help
  rip    Separate japanese text into words from CSV/TSV file
  rank   Show ranking of the word frequency
```

## Subcommands


### rip

`rip` command separate japanese text into words from `--input` file.

```sh
$ go-jp-text-ripper rip -h

Separate japanese text into words from CSV/TSV file

Options:

  -h, --help            display help information
  -c, --column          target column name in input file
      --columnn         target column index in input file (1st col=1)
  -i, --input          *input file path --input='/path/to/input.csv'
  -o, --output          output file path --output='./my_result.csv'
      --dic             custom dictionary path (mecab ipa dictionaly)
      --stopword        stop word list file path
      --show            print separated words to console
      --original        output original form of word
      --noun            output 'noun' type of word
      --verb            output 'verb' type of word
      --adjective       output 'adjective' type of word
      --neologd         use prefilter for neologd
      --progress[=30]   print current progress (sec)
      --min[=1]         minimum letter size for output
      --quote           columns to add double-quotes (separated by comma)
      --prefix          prefix name for new columns
  -r, --replace         replace from text column data to output result
      --debug           print debug result to console
      --dropempty       remove empty result from output
      --stoptop         use ranking from top as stopword
      --stoptopp        use ranking from top by percent as stopword (0.0 ~ 1.0)
      --stoplast        use ranking from last as stopword
      --stoplastp       use ranking from last by percent as stopword (0.0 ~ 1.0)
      --stopunique      use ranking stopword as unique per line
```

For example, if you want to separate words from the [example TSV file](example/aozora_bunko.tsv), try below command.

```sh
# chack the file contents
$ head -n 2 ./example/aozora_bunko.tsv

id	author	title	url	exerpt
1	夏目 漱石	吾輩は猫である	https://www.aozora.gr.jp/cards/000148/card789.html	一 吾輩は猫である。名前はまだ無い。 ...


# run rip command
$ go-jp-text-ripper rip \
    --input ./example/aozora_bunko.tsv \
    --column exerpt \
    --output ./output.tsv

[INFO]	[Run]	read and write lines...
[INFO]	[Run]	finish process

# check the results
$ head -n 2 ./output.tsv
id	author	title	url	exerpt	op_text	op_word_count	op_non_word_count	op_raw_char_count
1	夏目 漱石	吾輩は猫である	https://www.aozora.gr.jp/cards/000148/card789.html	一 吾輩は猫である。名前はまだ無い。...	一 吾輩 猫 名前 無い どこ 生れ 見当 つか 何 薄暗い し 所 ニャーニャー 泣い いた事 記憶 ...	562	719	2000
```


#### Advanced options

```sh
# `--columnn` sets column by index
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --show \
    --columnn 5

# `--dic` uses custom dictionary for kagome (https://github.com/ikawaha/kagome)
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --dic /opt/data/neologd.dic

# `--stopword` sets custom stopword file path and ignore the words
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --stopword ./stopwords.txt

# `--show` outputs the result on console
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt \
    --show

# `--original` uses original form (i.e. 原形) of the words for the results.
# in python code, use the word of `node.feature.split(",")[6]`
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --original

# if sets `--noun`, the results contains noun type of words.
# if sets `--verb`, the results contains verb type of words.
# if sets `--adjective`, the results contains adjective type of words.
# (default are 'noun', 'verb', 'adjective')
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --noun --verb  # in thie example, using only 'noun' and 'verb'

# `--neologd` uses the special prefilter for neologd to normalize text
# ref: https://github.com/evalphobia/go-jp-text-ripper/blob/master/prefilter/neologd.go
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --neologd

# `--progress` sets the interval in sec to show current progress
# default is '30' sec
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --progress 5

# `--min` sets the minimum letter size for the result
# if you set '2', then the result ignore one letter word (e.g. 'お', 'の', '犬', '猫', '嵐')
# default is '1'
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --min 3

# `--min` sets the minimum letter size for the result
# if you set '2', then the result ignore one letter word (e.g. 'お', 'の', '犬', '猫', '嵐')
# default is '1'
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --min 3

# `--prefix` sets the prefix for the new columns
# default is 'op_'
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --output ./output.tsv \
    --prefix n_

# `--replace` overwrite the target column by the result
# default is false and output the result on new column 'op_text'
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --min 3

# `--dropempty` removes the empty result row
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --dropempty

# `--stoptop`, `--stoptopp`, `--stoplast`, `--stoplastp` uses rank command result as a stopword
# `--stoptop` and `--stoptopp` uses the word with high frequency as a stopword
# `--stoplast` and `--stoplastp` uses the word with low frequency as a stopword
# if you use both of `--stoptop` and `--stoptopp` (or `--stoplast` and `--stoplastp`), then the filter condition stops when meets both.
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --stoptop 300
    --stoptopp 0.1  # whichever is bigger, 300 words or 10% words

# `--stopunique` is used with `--stop[top/last]` option
# this option count the frequency as one word per a row
$ go-jp-text-ripper rip --input ./example/aozora_bunko.tsv --column exerpt --show \
    --stoptop 300
    --stopunique
```

### rip

`rip` command separate japanese text into words from `--input` file.

```sh
$ go-jp-text-ripper rip -h

Separate japanese text into words from CSV/TSV file

Options:

  -h, --help            display help information
  -c, --column          target column name in input file
      --columnn         target column index in input file (1st col=1)
  -i, --input          *input file path --input='/path/to/input.csv'
  -o, --output          output file path --output='./my_result.csv'
      --dic             custom dictionary path (mecab ipa dictionaly)
      --stopword        stop word list file path
      --show            print separated words to console
      --original        output original form of word
      --noun            output 'noun' type of word
      --verb            output 'verb' type of word
      --adjective       output 'adjective' type of word
      --neologd         use prefilter for neologd
      --progress[=30]   print current progress (sec)
      --min[=1]         minimum letter size for output
      --quote           columns to add double-quotes (separated by comma)
      --prefix          prefix name for new columns
  -r, --replace         replace from text column data to output result
      --debug           print debug result to console
      --dropempty       remove empty result from output
      --stoptop         use ranking from top as stopword
      --stoptopp        use ranking from top by percent as stopword (0.0 ~ 1.0)
      --stoplast        use ranking from last as stopword
      --stoplastp       use ranking from last by percent as stopword (0.0 ~ 1.0)
      --stopunique      use ranking stopword as unique per line
```

For example, if you want to separate words from the [example TSV file](example/aozora_bunko.tsv), try below command.

```sh
# chack the file contents
$ head -n 2 ./example/aozora_bunko.tsv

id	author	title	url	exerpt
1	夏目 漱石	吾輩は猫である	https://www.aozora.gr.jp/cards/000148/card789.html	一 吾輩は猫である。名前はまだ無い。 ...


# run rip command
$ go-jp-text-ripper rip \
    --input ./example/aozora_bunko.tsv \
    --column exerpt \
    --output ./output.tsv

[INFO]	[Run]	read and write lines...
[INFO]	[Run]	finish process

# check the results
$ head -n 2 ./output.tsv
id	author	title	url	exerpt	op_text	op_word_count	op_non_word_count	op_raw_char_count
1	夏目 漱石	吾輩は猫である	https://www.aozora.gr.jp/cards/000148/card789.html	一 吾輩は猫である。名前はまだ無い。...	一 吾輩 猫 名前 無い どこ 生れ 見当 つか 何 薄暗い し 所 ニャーニャー 泣い いた事 記憶 ...	562	719	2000
```


## Custome Go App

Import `go-jp-text-ripper` and add plugins into `Config`.
You can add your custom plugins.

```go
package main

import (
	"github.com/evalphobia/go-jp-text-ripper/plugin"
	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// cli entry point
func main() {
	common := ripper.CommonConfig{}

	// prefilters to normalize raw text
	common.PreFilters = []*ripper.PreFilter{
			prefilter.Neologd,
	}

	// plugins
	common.Plugins = []*ripper.Plugin{
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
	common.PostFilters = []*ripper.PostFilter{
		postfilter.RatioJP,
		postfilter.RatioAlphaNum,
	}

	err := ripper.DoRip(ripper.RipConfig{
		CommonConfig:        common,
		DropEmpty:           true,
		StopWordTopNumber:   300,
	})
}
```

then, build and run!

# License

Apache License, Version 2.0
