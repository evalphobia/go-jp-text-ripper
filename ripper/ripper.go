package ripper

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/evalphobia/go-jp-text-ripper/normalizer"
	"github.com/evalphobia/go-jp-text-ripper/reader"
	"github.com/evalphobia/go-jp-text-ripper/tokenizer"
	"github.com/evalphobia/go-jp-text-ripper/writer"
)

const defaultPrefix = "op_"

// Prefix is output column prefix to add
var Prefix = defaultPrefix

// Ripper is struct for putting spaces between words
type Ripper struct {
	r           *reader.Reader
	inputHeader []string
	columnIndex int
	columnName  string

	w             *writer.Writer
	outputHeader  []string
	replaceColumn bool

	tok *tokenizer.Tokenizer
	nom *normalizer.Normalizer

	plugins []*Plugin

	quoteCols []string
	quoteIdx  []int

	ShowResult bool
	ShowDebug  bool
}

// New returns initialized Ripper
func New(col string) *Ripper {
	return &Ripper{
		columnName: col,
		tok:        tokenizer.New(),
		nom:        normalizer.Default,
	}
}

// NewFromFiles returns initialized Ripper
func NewFromFiles(in, out, col string) (*Ripper, error) {
	r := New(col)

	err := r.SetReaderFromFile(in)
	if err != nil {
		return nil, err
	}

	err = r.SetWriterFromFile(out)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// NewWithReaderFromFile returns initialized Ripper
func NewWithReaderFromFile(in, col string) (*Ripper, error) {
	r := New(col)

	err := r.SetReaderFromFile(in)
	if err != nil {
		return nil, err
	}

	r.w = writer.NewDummy()
	return r, nil
}

// SetReaderFromFile sets reader from file path
func (r *Ripper) SetReaderFromFile(path string) error {
	var err error
	r.r, err = reader.NewFromFile(path)
	return err
}

// SetWriterFromFile sets writer from file path
func (r *Ripper) SetWriterFromFile(path string) error {
	var err error
	r.w, err = writer.NewFromFile(path)
	return err
}

// SetDictionary sets dictinary
func (r *Ripper) SetDictionary(path string) error {
	return r.tok.SetDictionary(path)
}

// SetQuoteColumns sets normalizer
func (r *Ripper) SetQuoteColumns(cols []string) {
	c := make([]string, len(cols))
	for i, col := range cols {
		c[i] = strings.TrimSpace(col)
	}
	r.quoteCols = c
}

// SetNormalizer sets normalizer
func (r *Ripper) SetNormalizer(n *normalizer.Normalizer) {
	r.nom = n
}

// AddPlugin adds plugin
func (r *Ripper) AddPlugin(p *Plugin) {
	r.plugins = append(r.plugins, p)
}

// GetCurrentPosition return current pos
func (r *Ripper) GetCurrentPosition() int {
	return r.r.GetPosition()
}

// Close closes opened files
func (r *Ripper) Close() {
	r.r.Close()
	r.w.Close()
}

// ReadHeader reads header columns and check target column is existed or not
func (r *Ripper) ReadHeader(col string) error {
	header, err := r.r.Read()
	if err != nil {
		return err
	}

	hasColumn := false
	for idx, val := range header {
		if val == col {
			r.columnIndex = idx
			hasColumn = true
		}
		for _, q := range r.quoteCols {
			if val == q {
				r.quoteIdx = append(r.quoteIdx, idx)
				break
			}
		}
	}
	if !hasColumn {
		return fmt.Errorf("cannnot find column name in header: %s", col)
	}

	r.inputHeader = header
	return nil
}

// WriteHeader writes header columns
func (r *Ripper) WriteHeader() error {
	// read header if not read yet
	if len(r.inputHeader) == 0 {
		err := r.ReadHeader(r.columnName)
		if err != nil {
			return err
		}
	}

	// expand output header
	inHeader := r.inputHeader
	headerLen := len(inHeader)
	opHeader := make([]string, headerLen, headerLen+3)
	copy(opHeader, inHeader)

	// extra header name
	colText := Prefix + "text"
	colWordCount := Prefix + "word_count"
	colCharCount := Prefix + "raw_char_count"

	if !r.replaceColumn {
		opHeader = append(opHeader, colText)
	}

	extraHeaders := []string{colWordCount, colCharCount}
	for _, p := range r.plugins {
		extraHeaders = append(extraHeaders, p.Title)
	}
	r.outputHeader = append(opHeader, extraHeaders...)

	// write to file
	return r.w.Write(r.outputHeader)
}

// WriteHeaderWithReplace writes header columns and set as targe column is replaced
func (r *Ripper) WriteHeaderWithReplace() error {
	r.replaceColumn = true
	return r.WriteHeader()
}

// ReadAndWriteLines process each lines, read data, tokenize, and write it.
func (r *Ripper) ReadAndWriteLines() error {
	idx := r.columnIndex
	tok := r.tok
	for {
		line, err := r.r.Read()
		switch {
		case err == io.EOF:
			// end of file
			return nil
		case err != nil:
			return err
		}

		// tokenize text
		rawText := line[idx]
		normalizedText := r.nom.Normalize(rawText)
		tokens := tok.Tokenize(normalizedText)
		if err != nil {
			return err
		}

		if r.ShowDebug {
			showDebug(rawText, normalizedText, tokens)
		}

		// create result line
		words := tokens.GetWords()
		wordCount := strconv.Itoa(len(words))
		textLen := strconv.Itoa(utf8.RuneCountInString(rawText))
		wordLine := strings.Join(words, " ")
		if r.ShowResult {
			fmt.Println(wordLine)
		}

		var results []string
		if r.replaceColumn {
			line[idx] = wordLine
		} else {
			results = append(results, wordLine)
		}
		results = append(results, wordCount, textLen)

		// apply plugins
		for _, p := range r.plugins {
			pluginCount := p.Fn(rawText, normalizedText, tokens)
			results = append(results, pluginCount)
			if r.ShowDebug {
				fmt.Printf("%s: %s\n", p.Title, pluginCount)
			}
		}

		// quoting
		for _, i := range r.quoteIdx {
			line[i] = `"` + line[i] + `"`
		}

		// write result line
		results = append(line, results...)
		err = r.w.Write(results)
		if err != nil {
			return err
		}
	}
}

func showDebug(raw, nom string, list *tokenizer.TokenList) {
	const sep = "==============================\n"
	const sepMin = "------\n"
	fmt.Printf(sep)
	fmt.Printf("%s\n", raw)
	fmt.Printf(sepMin)
	fmt.Printf("%s\n", nom)
	fmt.Printf(sepMin)
	for _, t := range list.List {
		features := strings.Join(t.Token.Features(), ",")
		fmt.Printf("%s\t%v\n", t.Token.Surface, features)
	}
	fmt.Printf(sepMin)
}
