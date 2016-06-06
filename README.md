# go-jp-text-ripper

Separate long text into words and put spaces between ths words.


# Quick Usage

## command line
```sh
$ go get github.com/evalphobia/go-jp-text-ripper
$ go-jp-text-ripper -input ./raw_data.tsv -output ./result.tsv -column text_cell -replace
writeHeader...
processLines...
finish process

$ cat ./result.tsv
# // text_cell is separated into words!
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
- `-replace`: replace result
    - `false`: output result into added new column `sep_text`
    - `true`: output result into column `-column`
- `-show`: output separated words
- `-debug`: output word details

# License

Apache License, Version 2.0
