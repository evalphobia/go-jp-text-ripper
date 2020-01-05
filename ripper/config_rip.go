package ripper

import (
	"fmt"
)

const defaultPrefix = "op_"

// RipConfig contains options for 'rip' command.
type RipConfig struct {
	CommonConfig

	// columns to add double-quotes
	Quotes []string
	// replace the target column or not
	ReplaceText bool

	DropEmpty bool

	// use ranking from the top N as stopword
	StopWordTopNumber int
	// use ranking from the top by percent as stopword
	StopWordTopPercent float64 // 0.0~1.0
	// use ranking from the last N as stopword
	StopWordLastNumber int
	// use ranking from the last by percent as stopword
	StopWordLastPercent float64 // 0.0~1.0
	// use counting as one word if the same word exists in a line
	UseStopWordUnique bool
}

// Init initializes config.
func (c *RipConfig) Init() error {
	if c.Prefix == "" {
		c.Prefix = defaultPrefix
	}
	return c.CommonConfig.Init()
}

// Validate validates config.
func (c RipConfig) Validate() error {
	if err := c.CommonConfig.Validate(); err != nil {
		return err
	}

	if c.Output == "" && !c.ShowResult && !c.Debug {
		return fmt.Errorf("no output file\nSet -output <output file path> (or set -show option)")
	}
	return nil
}

// UseRankingForStopWord uses word frequency ranking as a stopword.
func (c RipConfig) UseRankingForStopWord() bool {
	switch {
	case c.StopWordTopNumber > 0,
		c.StopWordTopPercent > 0,
		c.StopWordLastNumber > 0,
		c.StopWordLastPercent > 0:
		return true
	}
	return false
}
