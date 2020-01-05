package main

import (
	"strings"

	"github.com/mkideal/cli"

	"github.com/evalphobia/go-jp-text-ripper/prefilter"
	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// rip command
type ripT struct {
	cli.Helper
	CommonOption
	Quote               string  `cli:"quote" usage:"columns to add double-quotes (separated by comma)"`
	Prefix              string  `cli:"prefix" usage:"prefix name for new columns"`
	ReplaceText         bool    `cli:"r,replace" usage:"replace from text column data to output result"`
	Debug               bool    `cli:"debug" usage:"print debug result to console"`
	DropEmpty           bool    `cli:"dropempty" usage:"remove empty result from output"`
	StopWordTopNumber   int     `cli:"stoptop" usage:"use ranking from top as stopword"`
	StopWordTopPercent  float64 `cli:"stoptopp" usage:"use ranking from top by percent as stopword (0.0 ~ 1.0)"`
	StopWordLastNumber  int     `cli:"stoplast" usage:"use ranking from last as stopword"`
	StopWordLastPercent float64 `cli:"stoplastp" usage:"use ranking from last by percent as stopword (0.0 ~ 1.0)"`
	UseStopWordUnique   bool    `cli:"stopunique" usage:"use ranking stopword as unique per line"`
}

var rip = &cli.Command{
	Name: "rip",
	Desc: "Separate japanese text into words from CSV/TSV file",
	Argv: func() interface{} { return new(ripT) },
	Fn:   execRip,
}

func execRip(ctx *cli.Context) error {
	argv := ctx.Argv().(*ripT)

	common := ripper.CommonConfig{
		Column:           argv.Column,
		ColumnNumber:     argv.ColumnNumber,
		Input:            argv.Input,
		Output:           argv.Output,
		Dictionary:       argv.Dictionary,
		StopWordPath:     argv.StopWord,
		ShowResult:       argv.ShowResult,
		UseOriginalForm:  argv.UseOriginalForm,
		UseNoun:          argv.UseNoun,
		UseVerb:          argv.UseVerb,
		UseAdjective:     argv.UseAdjective,
		UseNeologd:       argv.UseNeologd,
		ProgressInterval: argv.ProgressInterval,
		MinLetterSize:    argv.MinLetterSize,
		Version:          version,
		Revision:         revision,
		Prefix:           argv.Prefix,
		Debug:            argv.Debug,
	}
	if common.UseNeologd {
		common.PreFilters = append(common.PreFilters, prefilter.Neologd)
	}
	return ripper.DoRip(ripper.RipConfig{
		CommonConfig:        common,
		ReplaceText:         argv.ReplaceText,
		Quotes:              strings.Split(argv.Quote, ","),
		DropEmpty:           argv.DropEmpty,
		StopWordTopNumber:   argv.StopWordTopNumber,
		StopWordTopPercent:  argv.StopWordTopPercent,
		StopWordLastNumber:  argv.StopWordLastNumber,
		StopWordLastPercent: argv.StopWordLastPercent,
		UseStopWordUnique:   argv.UseStopWordUnique,
	})
}
