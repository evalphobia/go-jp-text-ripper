package main

import (
	"flag"
	"fmt"

	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// flags
var (
	input   = ""
	output  = ""
	column  = ""
	dic     = ""
	replace = false
	show    = false
	debug   = false
)

// cli entry point
func main() {
	err := initFlags()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var r *ripper.Ripper
	switch {
	case output == "":
		// read only
		r, err = ripper.NewWithReaderFromFile(input, column)
	default:
		r, err = ripper.NewFromFiles(input, output, column)
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer r.Close()

	// set output options
	if show {
		r.ShowResult = true
	}
	if debug {
		r.ShowDebug = true
	}

	// set original dictionary
	if dic != "" {
		err = r.SetDictionary(dic)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	switch {
	case replace:
		err = r.WriteHeaderWithReplace()
	default:
		err = r.WriteHeader()
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("read and write lines...")
	err = r.ReadAndWriteLines()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("finish process")
}

func initFlags() error {
	err := parseFlags()
	if err != nil {
		return err
	}

	return checkFlags()
}

func parseFlags() error {
	flag.StringVar(&input, "input", "", "read file")
	flag.StringVar(&output, "output", "", "output file")
	flag.StringVar(&column, "column", "", "target column name")
	flag.StringVar(&dic, "dic", "", "custom dictionaly path(ipa dictionaly)")
	flag.BoolVar(&replace, "replace", false, "replace text column")
	flag.BoolVar(&show, "show", false, "print separated words to console")
	flag.BoolVar(&debug, "debug", false, "print debug result to console")

	flag.Parse()
	return nil
}

func checkFlags() error {
	switch {
	case input == "":
		return fmt.Errorf("no input file\nuse -csv <input file path>\n")
	case output == "" && !show && !debug:
		return fmt.Errorf("no output file\nuse -output <output file path>\n")
	case column == "":
		return fmt.Errorf("no column name\nuse -column <column name>\n")
	}

	return nil
}
