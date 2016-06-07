package ripper

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/evalphobia/go-jp-text-ripper/reader"
	"github.com/evalphobia/go-jp-text-ripper/tokenizer"
	"github.com/evalphobia/go-jp-text-ripper/writer"
)

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
	return r.tok.SetDictinary(path)
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

	extraHeaders := []string{"op_word_count", "raw_char_count"}
	if !r.replaceColumn {
		extraHeaders = append([]string{"sep_text"}, extraHeaders...)
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
		text := line[idx]
		tokens := tok.Tokenize(text)
		if err != nil {
			return err
		}

		if r.ShowDebug {
			showDebug(text, tokens)
		}

		words := tokens.GetWords()
		wordCount := strconv.Itoa(len(words))
		textLen := strconv.Itoa(len(text))

		wordLine := strings.Join(words, " ")
		if r.ShowResult {
			fmt.Println(wordLine)
		}

		// create new line
		var result []string
		if r.replaceColumn {
			line[idx] = wordLine
			result = append(line, wordCount, textLen)
		} else {
			result = append(line, wordLine, wordCount, textLen)
		}

		err = r.w.Write(result)
		if err != nil {
			return err
		}
	}
}

func showDebug(text string, list *tokenizer.TokenList) {
	const sep = "==============================\n"
	fmt.Printf(sep)
	fmt.Printf("%s\n\n", text)
	for _, t := range list.List {
		features := strings.Join(t.Token.Features(), ",")
		fmt.Printf("%s\t%v\n", t.Token.Surface, features)
	}
	fmt.Printf(sep)
}
