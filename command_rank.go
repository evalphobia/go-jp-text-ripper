package main

import (
	"github.com/mkideal/cli"

	"github.com/evalphobia/go-jp-text-ripper/prefilter"
	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// rank command
type rankT struct {
	cli.Helper
	CommonOption
	TopNumber   int     `cli:"top" usage:"rank from top by count"`
	TopPercent  float64 `cli:"topp" usage:"rank from top by percent (0.0 ~ 1.0)"`
	LastNumber  int     `cli:"last" usage:"rank from last by count"`
	LastPercent float64 `cli:"lastp" usage:"rank from last by percent (0.0 ~ 1.0)"`
	UseUnique   bool    `cli:"u,unique" usage:"count as one word if the same word exists in a line"`
}

var rank = &cli.Command{
	Name: "rank",
	Desc: "Show ranking of the word frequency",
	Argv: func() interface{} { return new(rankT) },
	Fn:   execRank,
}

func execRank(ctx *cli.Context) error {
	argv := ctx.Argv().(*rankT)

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
	}
	if common.UseNeologd {
		common.PreFilters = append(common.PreFilters, prefilter.Neologd)
	}
	return ripper.DoRank(ripper.RankConfig{
		CommonConfig: common,
		TopNumber:    argv.TopNumber,
		TopPercent:   argv.TopPercent,
		LastNumber:   argv.LastNumber,
		LastPercent:  argv.LastPercent,
		UseUnique:    argv.UseUnique,
	})
}
