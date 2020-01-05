package ripper

import (
	"fmt"
)

const (
	defaultTopNumber = 100
)

// RankConfig contains options for 'rank' command.
type RankConfig struct {
	CommonConfig

	// show ranking from the top N
	TopNumber int
	// show ranking from the top by percent
	TopPercent float64 // 0.0~1.0
	// show ranking from the last N
	LastNumber int
	// show ranking from the last by percent
	LastPercent float64 // 0.0~1.0

	// count as one word if the same word exists in a line.
	UseUnique bool
}

// Init initializes config.
func (c *RankConfig) Init() error {
	switch {
	case c.TopNumber > 0,
		c.TopPercent > 0,
		c.LastNumber > 0,
		c.LastPercent > 0:
		// pass
	default:
		c.TopNumber = defaultTopNumber
	}

	return c.CommonConfig.Init()
}

// Validate validates config.
func (c RankConfig) Validate() error {
	if err := c.CommonConfig.Validate(); err != nil {
		return err
	}

	if c.Output == "" && !c.ShowResult {
		return fmt.Errorf("no output file\nSet -output <output file path> (or set -show option)")
	}
	return nil
}
