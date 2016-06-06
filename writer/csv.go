package writer

import (
	"encoding/csv"
	"os"
	"path"
)

// Writer writes the output to file
type Writer struct {
	fp *os.File
	w  writer
}

// NewFromFile returns initialized Writer for file
func NewFromFile(filepath string) (*Writer, error) {
	fp, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	var w writer
	switch ext := path.Ext(filepath); ext {
	case ".tsv":
		w = newTSVWriter(fp)
	default:
		w = newCSVWriter(fp)
	}

	return &Writer{
		fp: fp,
		w:  w,
	}, nil
}

// NewDummy returns initialized Writer with dummy writer
func NewDummy() *Writer {
	return &Writer{
		fp: nil,
		w:  newDummyWriter(),
	}
}

// Write writes a line into file
func (w *Writer) Write(line []string) error {
	err := w.w.Write(line)
	if err != nil {
		return err
	}

	w.w.Flush()
	return nil
}

// Close closes file
func (w *Writer) Close() {
	if w.fp == nil {
		return
	}
	w.fp.Close()
}

// writer is interface of actual writes line into files
type writer interface {
	Write([]string) error
	Flush()
}

func newCSVWriter(fp *os.File) writer {
	w := csv.NewWriter(fp)
	return w
}

func newTSVWriter(fp *os.File) writer {
	w := csv.NewWriter(fp)
	w.Comma = '\t'
	return w
}
