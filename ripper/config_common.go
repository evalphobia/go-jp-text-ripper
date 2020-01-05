package ripper

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/evalphobia/go-jp-text-ripper/log"
	"github.com/evalphobia/go-jp-text-ripper/tokenizer"
)

// CommonConfig contains common options for subcommands.
type CommonConfig struct {
	// input file path
	Input string
	// output file path
	Output string
	// target column name
	Column string
	// target column index number (first=1)
	ColumnNumber int
	// custome dictionary for ikawaha/kagome
	Dictionary string
	// print separated words to console
	ShowResult bool
	// intervals to print current progress (sec)
	ProgressInterval int
	// print debug result to console
	Debug bool
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

	// Version info
	Version  string
	Revision string

	Logger log.Logger
}

// Init initializes config.
func (c *CommonConfig) Init() error {
	if c.Logger == nil {
		c.Logger = &log.StdLogger{}
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

// Validate validates config.
func (c CommonConfig) Validate() error {
	switch {
	case c.Column == "" && c.ColumnNumber == 0:
		return fmt.Errorf("no target column\nSet -column <column name> (or -columnn <column index>)")
	case c.Input == "":
		return fmt.Errorf("no input file\nSet -input <input file path>")
		// case c.Output == "" && !c.ShowResult && !c.Debug:
		// return fmt.Errorf("no output file\nSet -output <output file path> (or set -show option)\n")
	}

	return nil
}

// GetPosList returns 'the parts of speech' for tokenizer.
func (c CommonConfig) GetPosList() []string {
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
