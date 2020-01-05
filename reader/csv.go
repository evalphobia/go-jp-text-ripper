package reader

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
)

// Reader reads file
type Reader struct {
	fp       *os.File
	r        reader
	position int
}

// NewFromFile returns initialized Reader for file
func NewFromFile(filepath string) (*Reader, error) {
	/* #nosec G304 */
	fp, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	var r reader
	switch ext := path.Ext(filepath); ext {
	case ".csv":
		r = newCSVReader(fp)
	case ".tsv":
		r = newTSVReader(fp)
	default:
		return nil, fmt.Errorf("non supported file format: %s", ext)
	}

	return &Reader{
		fp: fp,
		r:  r,
	}, nil
}

// Read returns []string and count up current position
func (r *Reader) Read() ([]string, error) {
	line, err := r.r.Read()
	if err != nil {
		return nil, err
	}

	r.position++
	return line, nil
}

// Close closes file
func (r *Reader) Close() error {
	return r.fp.Close()
}

// GetPosition returns position(read line number)
func (r *Reader) GetPosition() int {
	return r.position
}

// reader is interface of actual reads line from files
type reader interface {
	Read() ([]string, error)
}

func newCSVReader(fp *os.File) reader {
	r := csv.NewReader(fp)
	r.FieldsPerRecord = -1
	return r
}

func newTSVReader(fp *os.File) reader {
	r := csv.NewReader(fp)
	r.Comma = '\t'
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	return r
}
