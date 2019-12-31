package ripper

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/evalphobia/go-jp-text-ripper/log"
	"github.com/evalphobia/go-jp-text-ripper/reader"
	"github.com/evalphobia/go-jp-text-ripper/tokenizer"
	"github.com/evalphobia/go-jp-text-ripper/writer"
)

// Ripper is struct for putting spaces between words
type Ripper struct {
	r           *reader.Reader
	inputHeader []string
	columnIndex int

	w            *writer.Writer
	outputHeader []string

	tok         *tokenizer.Tokenizer
	preFilters  []*PreFilter
	plugins     []*Plugin
	postFilters []*PostFilter

	quoteCols []string
	quoteIdx  []int

	Config Config
}

// New returns initialized Ripper.
func New(c Config) (*Ripper, error) {
	r := &Ripper{
		tok: tokenizer.New(tokenizer.Config{
			WordPosList:     c.GetPosList(),
			StopWordList:    c.StopWords,
			MinLetterSize:   c.MinLetterSize,
			UseOriginalForm: c.UseOriginalForm,
		}),
	}

	if err := r.SetReaderFromFile(c.Input); err != nil {
		return nil, err
	}
	switch {
	case c.Output == "":
		r.w = writer.NewDummy()
	default:
		if err := r.SetWriterFromFile(c.Output); err != nil {
			return nil, err
		}
	}

	r.Config = c
	if len(c.Quotes) != 0 {
		r.SetQuoteColumns(c.Quotes)
	}

	// set original dictionary
	if c.Dictionary != "" {
		if err := r.SetDictionary(c.Dictionary); err != nil {
			r.Close()
			return nil, err
		}
	}

	r.AddPreFilters(c.PreFilters...)
	r.AddPlugins(c.Plugins...)
	r.AddPostFilters(c.PostFilters...)
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

// AddPreFilter adds pre filter.
func (r *Ripper) AddPreFilters(p ...*PreFilter) {
	r.preFilters = append(r.preFilters, p...)
}

// AddPlugin adds plugin.
func (r *Ripper) AddPlugins(p ...*Plugin) {
	r.plugins = append(r.plugins, p...)
}

// AddPostFilter adds post filter.
func (r *Ripper) AddPostFilters(p ...*PostFilter) {
	r.postFilters = append(r.postFilters, p...)
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
		return fmt.Errorf("cannnot find column name in header: col:[%s] headers:[%+v]", col, header)
	}

	r.inputHeader = header
	return nil
}

// WriteHeader writes header columns
func (r *Ripper) WriteHeader() error {
	c := r.Config

	// read header if not read yet
	if len(r.inputHeader) == 0 {
		err := r.ReadHeader(c.Column)
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
	colText := c.Prefix + "text"
	colWordCount := c.Prefix + "word_count"
	colNonWordCount := c.Prefix + "non_word_count"
	colCharCount := c.Prefix + "raw_char_count"

	if !r.Config.ReplaceText {
		opHeader = append(opHeader, colText)
	}

	extraHeaders := []string{colWordCount, colNonWordCount, colCharCount}
	for _, p := range r.plugins {
		extraHeaders = append(extraHeaders, c.Prefix+p.Title)
	}
	for _, p := range r.postFilters {
		extraHeaders = append(extraHeaders, c.Prefix+p.Title)
	}

	r.outputHeader = append(opHeader, extraHeaders...)

	// write to file
	return r.w.Write(r.outputHeader)
}

// ReadAndWriteLines process each lines, read data, tokenize, and write it.
func (r *Ripper) ReadAndWriteLines() error {
	c := r.Config
	logger := c.Logger
	idx := r.columnIndex
	tok := r.tok

	lastLineNo := 1
	lastLineText := ""
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		logger.Errorf("ReadAndWriteLines", "unknown error occurred on Line:[%d] Text:[%s]\n", lastLineNo, lastLineText)
	}()

	for {
		lastLineNo++
		line, err := r.r.Read()
		switch {
		case err == io.EOF:
			// end of file
			return nil
		case err != nil:
			logger.Errorf("ReadAndWriteLines", "r.r.Read() err:[%s]\n", err.Error())
			return err
		}

		text := &TextData{}

		// tokenize text
		lastLineText = line[idx]
		text.raw = line[idx]
		text.normalized = r.applyPreFilters(text.raw)
		text.words, text.nonWords = tok.Tokenize(text.normalized)
		if err != nil {
			return err
		}

		if c.Debug {
			showDebug(logger, text)
		}

		// create result line
		words := text.words.GetWords()
		wordCount := strconv.Itoa(len(words))
		nonWordCount := strconv.Itoa(len(text.nonWords.GetWords()))
		textLen := strconv.Itoa(utf8.RuneCountInString(text.raw))
		wordLine := strings.Join(words, " ")
		if c.ShowResult {
			logger.Infof("ReadAndWriteLines", wordLine)
		}
		if c.DropEmpty && wordLine == "" {
			continue
		}

		var results []string
		if c.ReplaceText {
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
			logger.Errorf("ReadAndWriteLines", "r.w.Write() err:[%s]\n", err.Error())
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
	c := r.Config
	logger := c.Logger

	for _, p := range r.plugins {
		fnResult := p.Fn(text)
		results = append(results, fnResult)
		if c.Debug {
			logger.Debugf("applyPlugins", "%s: %s\n", p.Title, fnResult)
		}
	}
	return results
}

// applyPostFilters runs postfilters function and adds the result
func (r *Ripper) applyPostFilters(results, line []string) []string {
	if len(r.postFilters) == 0 {
		return results
	}
	c := r.Config
	logger := c.Logger

	data := make(map[string]string)
	header := r.outputHeader
	for i, val := range append(line, results...) {
		title := strings.TrimPrefix(header[i], c.Prefix)
		data[title] = val
	}

	for _, p := range r.postFilters {
		fnResult := p.Fn(data)
		results = append(results, fnResult)
		if c.Debug {
			logger.Debugf("applyPostFilters", "%s: %s\n", p.Title, fnResult)
		}
	}
	return results
}

func showDebug(logger log.Logger, text *TextData) {
	const sep = "=============================="
	const sepMin = "------"
	data := make([]string, 0, 1024)
	data = append(data, sep)
	data = append(data, text.raw)
	data = append(data, sepMin)
	data = append(data, text.normalized)
	data = append(data, fmt.Sprintf("%s words: %d", sepMin, len(text.words.List)))
	for _, t := range text.words.List {
		features := strings.Join(t.Token.Features(), ",")
		data = append(data, fmt.Sprintf("%s\t%v", t.Token.Surface, features))
	}
	data = append(data, fmt.Sprintf("%s non-words: %d\n", sepMin, len(text.words.List)))
	for _, t := range text.nonWords.List {
		features := strings.Join(t.Token.Features(), ",")
		data = append(data, fmt.Sprintf("%s\t%v", t.Token.Surface, features))
	}
	data = append(data, sepMin)
	logger.Debugf("", strings.Join(data, "\n"))
}
