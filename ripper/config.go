package ripper

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/evalphobia/go-jp-text-ripper/log"
	"github.com/evalphobia/go-jp-text-ripper/tokenizer"
)

const defaultPrefix = "op_"

type Config struct {
	// input file path
	Input string
	// output file path
	Output string
	// target column name
	Column string
	// custome dictionary for ikawaha/kagome
	Dictionary string
	// columns to add double-quotes
	Quotes []string
	// replace the target column or not
	ReplaceText bool
	// print separated words to console
	ShowResult bool
	// print debug result to console
	Debug bool
	// intervals to print current progress (sec)
	ProgressInterval int
	// Prefix is output column prefix to add
	Prefix string

	// Filter and Plugins
	PreFilters  []*PreFilter
	Plugins     []*Plugin
	PostFilters []*PostFilter
	UseNeologd  bool

	// Tokenizer settings:
	MinLetterSize   int
	StopWordPath    string
	StopWords       []string
	UseOriginalForm bool
	UseNoun         bool
	UseVerb         bool
	UseAdjective    bool
	DropEmpty       bool

	// Version info
	Version  string
	Revision string

	Logger log.Logger
}

func (c *Config) Init() error {
	if flag.Parsed() {
		return nil
	}

	flag.StringVar(&c.Input, "input", "", "read file")
	flag.StringVar(&c.Output, "output", "", "output file")
	flag.StringVar(&c.Column, "column", "", "target column name")
	flag.StringVar(&c.Dictionary, "dic", "", "custom dictionaly path(ipa dictionaly)")
	flag.StringVar(&c.StopWordPath, "stopword", "", "stop word list file path")
	flag.BoolVar(&c.ReplaceText, "replace", false, "replace text column")
	flag.BoolVar(&c.ShowResult, "show", false, "print separated words to console")
	flag.BoolVar(&c.Debug, "debug", false, "print debug result to console")
	flag.BoolVar(&c.UseOriginalForm, "original", false, "output original form of word")
	flag.BoolVar(&c.UseNoun, "noun", false, "keep 'noun' as a word")
	flag.BoolVar(&c.UseVerb, "verb", false, "keep 'verb' as a word")
	flag.BoolVar(&c.UseAdjective, "adjective", false, "keep 'adjective' as a word")
	flag.BoolVar(&c.UseNeologd, "neologd", false, "use prefilter for neologd")
	flag.BoolVar(&c.DropEmpty, "dropempty", false, "remove empty output word line")
	flag.IntVar(&c.ProgressInterval, "progress", 30, "print current progress (sec)")
	flag.IntVar(&c.MinLetterSize, "min", 1, "minimum letter size to keep as a word")

	var quote string
	flag.StringVar(&quote, "quote", "", "columns to add double-quotes (separated by comma)")
	flag.Parse()

	if quote != "" {
		c.Quotes = strings.Split(quote, ",")
	}
	if c.Logger == nil {
		c.Logger = &log.StdLogger{}
	}
	if c.Prefix == "" {
		c.Prefix = defaultPrefix
	}
	switch {
	case c.StopWordPath != "":
		words, err := getWordsFromPath(c.StopWordPath)
		if err != nil {
			return err
		}
		c.StopWords = append(c.StopWords, words...)
	}

	return nil
}

func (c Config) Validate() error {
	switch {
	case c.Input == "":
		return fmt.Errorf("no input file\nuse -input <input file path>\n")
	case c.Output == "" && !c.ShowResult && !c.Debug:
		return fmt.Errorf("no output file\nuse -output <output file path>\n")
	case c.Column == "":
		return fmt.Errorf("no column name\nuse -column <column name>\n")
	}

	return nil
}

// return 'the parts of speech' for tokenizer.
func (c Config) GetPosList() []string {
	var pos []string
	if c.UseNoun {
		pos = append(pos, tokenizer.PosNoun)
	}
	if c.UseVerb {
		pos = append(pos, tokenizer.PosVerb)
	}
	if c.UseAdjective {
		pos = append(pos, tokenizer.PosAdjective)
	}
	return pos
}

func getWordsFromPath(path string) ([]string, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0, 1024)
	r := bufio.NewReaderSize(fp, 4096)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}
