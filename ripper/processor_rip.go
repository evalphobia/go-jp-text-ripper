package ripper

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/evalphobia/go-jp-text-ripper/log"
)

// DoRip creates *RipProcessor from config and run it.
func DoRip(conf RipConfig) error {
	if err := conf.Init(); err != nil {
		return err
	}
	if err := conf.Validate(); err != nil {
		return err
	}

	conf.Logger.Infof("DoRip", "version:[%s] rev:[%s]", conf.Version, conf.Revision)
	r, err := NewRipProcessor(conf)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := r.WriteHeader(); err != nil {
		r.Close()
		return err
	}

	return r.DoWithProgress()
}

// RipProcessor is struct for putting spaces between words.
type RipProcessor struct {
	*CommonProcessor
	Config    RipConfig
	quoteCols []string
	quoteIdx  []int
}

// NewRipProcessor returns initialized RipProcessor.
func NewRipProcessor(c RipConfig) (*RipProcessor, error) {
	common, err := NewCommonProcessor(c.CommonConfig)
	if err != nil {
		return nil, err
	}

	r := &RipProcessor{
		CommonProcessor: common,
		Config:          c,
	}
	if len(c.Quotes) != 0 {
		r.SetQuoteColumns(c.Quotes)
	}
	return r, nil
}

// SetQuoteColumns sets normalizer
func (r *RipProcessor) SetQuoteColumns(cols []string) {
	c := make([]string, len(cols))
	for i, col := range cols {
		c[i] = strings.TrimSpace(col)
	}
	r.quoteCols = c
}

// ReadHeader reads header columns and sets target column.
func (r *RipProcessor) ReadHeader() error {
	c := r.Config
	switch {
	case c.ColumnNumber > 0:
		return r.CommonProcessor.ReadHeaderWithIndex(c.ColumnNumber - 1)
	default:
		return r.readHeaderByName(c.Column)
	}
}

// readHeaderByName reads header columns and check target column is existed or not.
func (r *RipProcessor) readHeaderByName(col string) error {
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
func (r *RipProcessor) WriteHeader() error {
	c := r.Config

	// read header if not read yet
	if len(r.inputHeader) == 0 {
		err := r.ReadHeader()
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

// DoWithProgress processes with showing progress.
func (r *RipProcessor) DoWithProgress() error {
	r.ShowProgress()

	conf := r.Config
	logger := conf.Logger
	logger.Infof("Run", "read and write lines...")

	err := r.Do()
	if err != nil {
		logger.Errorf("Run", "error on r.Process() err:[%s]", err.Error())
		return err
	}

	logger.Infof("Run", "finish process")
	return nil
}

// Do processes each lines, read data, tokenize, and write it.
func (r *RipProcessor) Do() error {
	defer r.Close()
	c := r.Config
	logger := c.Logger
	if c.UseRankingForStopWord() {
		rank, err := r.doGetRankStopWord()
		if err != nil {
			return err
		}
		r.tok.AddStopWords(rank.GetTopWords()...)
		r.tok.AddStopWords(rank.GetLastWords()...)
	}

	idx := r.columnIndex
	tok := r.tok

	lastLineNo := 1
	lastLineText := ""
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		logger.Errorf("Do", "unknown error occurred on Line:[%d] Text:[%s]\n", lastLineNo, lastLineText)
	}()

	for {
		lastLineNo++
		line, err := r.r.Read()
		switch {
		case err == io.EOF:
			// end of file
			return nil
		case err != nil:
			logger.Errorf("Do", "r.r.Read() err:[%s]\n", err.Error())
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
			logger.Infof("Do", wordLine)
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
			logger.Errorf("Do", "r.w.Write() err:[%s]\n", err.Error())
			return err
		}
	}
}

// doGetRankStopWord gets word frequency for the stop words.
func (r *RipProcessor) doGetRankStopWord() (RankResult, error) {
	c := r.Config
	rp, err := NewRankProcessor(RankConfig{
		CommonConfig: c.CommonConfig,
		TopNumber:    c.StopWordTopNumber,
		TopPercent:   c.StopWordTopPercent,
		LastNumber:   c.StopWordLastNumber,
		LastPercent:  c.StopWordLastPercent,
		UseUnique:    c.UseStopWordUnique,
	})
	if err != nil {
		return RankResult{}, err
	}

	return rp.GetRank()
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
