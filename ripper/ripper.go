package ripper

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

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

	tok         *tokenizer.Tokenizer
	preFilters  []*PreFilter
	plugins     []*Plugin
	postFilters []*PostFilter

	quoteCols []string
	quoteIdx  []int

	// console output flags
	ShowResult bool
	ShowDebug  bool
}

// New returns initialized Ripper
func New(col string) *Ripper {
	return &Ripper{
		columnName: col,
		tok:        tokenizer.New(),
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

// AddPreFilter adds pre filter
func (r *Ripper) AddPreFilter(p *PreFilter) {
	r.preFilters = append(r.preFilters, p)
}

// AddPlugin adds plugin
func (r *Ripper) AddPlugin(p *Plugin) {
	r.plugins = append(r.plugins, p)
}

// AddPostFilter adds post filter
func (r *Ripper) AddPostFilter(p *PostFilter) {
	r.postFilters = append(r.postFilters, p)
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
	colNonWordCount := Prefix + "non_word_count"
	colCharCount := Prefix + "raw_char_count"

	if !r.replaceColumn {
		opHeader = append(opHeader, colText)
	}

	extraHeaders := []string{colWordCount, colNonWordCount, colCharCount}
	for _, p := range r.plugins {
		extraHeaders = append(extraHeaders, Prefix+p.Title)
	}
	for _, p := range r.postFilters {
		extraHeaders = append(extraHeaders, Prefix+p.Title)
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

		text := &TextData{}

		// tokenize text
		text.raw = line[idx]
		text.normalized = r.applyPreFilters(text.raw)
		text.words, text.nonWords = tok.Tokenize(text.normalized)
		if err != nil {
			return err
		}

		if r.ShowDebug {
			showDebug(text)
		}

		// create result line
		words := text.words.GetWords()
		wordCount := strconv.Itoa(len(words))
		nonWordCount := strconv.Itoa(len(text.nonWords.GetWords()))
		textLen := strconv.Itoa(utf8.RuneCountInString(text.raw))
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
		results = append(results, wordCount, nonWordCount, textLen)

		results = r.applyPlugins(results, text)
		results = r.applyPostFilters(results, line)

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

// applyPreFilters runs prefilters function and return normalized text
func (r *Ripper) applyPreFilters(text string) string {
	for _, p := range r.preFilters {
		text = p.Fn(text)
	}
	return text
}

// applyPlugins runs plugins function and adds result
func (r *Ripper) applyPlugins(results []string, text *TextData) []string {
	for _, p := range r.plugins {
		fnResult := p.Fn(text)
		results = append(results, fnResult)
		if r.ShowDebug {
			fmt.Printf("%s: %s\n", p.Title, fnResult)
		}
	}
	return results
}

// applyPostFilters runs postfilters function and adds the result
func (r *Ripper) applyPostFilters(results, line []string) []string {
	if len(r.postFilters) == 0 {
		return results
	}

	data := make(map[string]string)
	header := r.outputHeader
	for i, val := range append(line, results...) {
		title := strings.TrimPrefix(header[i], Prefix)
		data[title] = val
	}

	for _, p := range r.postFilters {
		fnResult := p.Fn(data)
		results = append(results, fnResult)
		if r.ShowDebug {
			fmt.Printf("%s: %s\n", p.Title, fnResult)
		}
	}
	return results
}

func showDebug(text *TextData) {
	const sep = "==============================\n"
	const sepMin = "------\n"
	fmt.Printf(sep)
	fmt.Printf("%s\n", text.raw)
	fmt.Printf(sepMin)
	fmt.Printf("%s\n", text.normalized)
	fmt.Printf("%s words: %d\n", sepMin, len(text.words.List))
	for _, t := range text.words.List {
		features := strings.Join(t.Token.Features(), ",")
		fmt.Printf("%s\t%v\n", t.Token.Surface, features)
	}
	fmt.Printf("%s non-words: %d\n", sepMin, len(text.words.List))
	for _, t := range text.nonWords.List {
		features := strings.Join(t.Token.Features(), ",")
		fmt.Printf("%s\t%v\n", t.Token.Surface, features)
	}
	fmt.Printf(sepMin)
}
