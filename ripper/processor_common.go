package ripper

import (
	"strings"
	"time"

	"github.com/evalphobia/go-jp-text-ripper/reader"
	"github.com/evalphobia/go-jp-text-ripper/tokenizer"
	"github.com/evalphobia/go-jp-text-ripper/writer"
)

// CommonProcessor is common struct for processing.
type CommonProcessor struct {
	r           *reader.Reader
	inputHeader []string
	columnIndex int

	w            *writer.Writer
	outputHeader []string

	tok         *tokenizer.Tokenizer
	preFilters  []*PreFilter
	plugins     []*Plugin
	postFilters []*PostFilter

	Config CommonConfig
}

// NewCommonProcessor returns initialized CommonProcessor.
func NewCommonProcessor(c CommonConfig) (*CommonProcessor, error) {
	r := &CommonProcessor{
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
func (r *CommonProcessor) SetReaderFromFile(path string) error {
	var err error
	r.r, err = reader.NewFromFile(path)
	return err
}

// SetWriterFromFile sets writer from file path
func (r *CommonProcessor) SetWriterFromFile(path string) error {
	var err error
	r.w, err = writer.NewFromFile(path)
	return err
}

// SetDictionary sets dictinary
func (r *CommonProcessor) SetDictionary(path string) error {
	return r.tok.SetDictionary(path)
}

// SetColumnIndex sets index of column (first=0).
func (r *CommonProcessor) SetColumnIndex(idx int) {
	r.columnIndex = idx
}

// AddPreFilters adds pre filter.
func (r *CommonProcessor) AddPreFilters(p ...*PreFilter) {
	r.preFilters = append(r.preFilters, p...)
}

// AddPlugins adds plugin.
func (r *CommonProcessor) AddPlugins(p ...*Plugin) {
	r.plugins = append(r.plugins, p...)
}

// AddPostFilters adds post filter.
func (r *CommonProcessor) AddPostFilters(p ...*PostFilter) {
	r.postFilters = append(r.postFilters, p...)
}

// GetCurrentPosition return current pos
func (r *CommonProcessor) GetCurrentPosition() int {
	return r.r.GetPosition()
}

// Close closes opened files
func (r *CommonProcessor) Close() {
	r.r.Close()
	r.w.Close()
}

// ReadHeader reads column of header from input file.
func (r *CommonProcessor) ReadHeader() error {
	header, err := r.r.Read()
	if err != nil {
		return err
	}
	r.inputHeader = header
	return nil
}

// ReadHeaderWithIndex reads header columns and sets target column by index.
func (r *CommonProcessor) ReadHeaderWithIndex(idx int) error {
	err := r.ReadHeader()
	if err != nil {
		return err
	}

	r.SetColumnIndex(idx)
	return nil
}

// applyPreFilters runs prefilters function and return normalized text
func (r *CommonProcessor) applyPreFilters(text string) string {
	for _, p := range r.preFilters {
		text = p.Fn(text)
	}
	return text
}

// applyPlugins runs plugins function and adds result
func (r *CommonProcessor) applyPlugins(results []string, text *TextData) []string {
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
func (r *CommonProcessor) applyPostFilters(results, line []string) []string {
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

// ShowProgress outputs current progress on background.
func (r *CommonProcessor) ShowProgress() {
	conf := r.Config
	logger := conf.Logger
	go func() {
		interval := conf.ProgressInterval
		tick := time.Tick(time.Duration(conf.ProgressInterval) * time.Second)
		prev := 0
		for {
			for range tick {
				cur := r.GetCurrentPosition()
				logger.Infof("progress", "line: %d, tps: %d\n", cur, (cur-prev)/interval)
				prev = cur
			}
		}
	}()
}
