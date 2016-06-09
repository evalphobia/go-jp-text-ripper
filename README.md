# go-jp-text-ripper

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
