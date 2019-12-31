package ripper

import (
	"fmt"
	"time"
)

// AutoRun creates *Ripper from CLI flags and run it
func AutoRun(conf Config) {
	if err := conf.Init(); err != nil {
		fmt.Printf("error on conf.Init(). err:[%s]", err.Error())
		return
	}
	if err := conf.Validate(); err != nil {
		conf.Logger.Errorf("AutoRun", "error on [conf.Validate()]. err:[%s]", err.Error())
		return
	}

	conf.Logger.Infof("AutoRun", "version:[%s] rev:[%s]", conf.Version, conf.Revision)
	r, err := newDefaultRipper(conf)
	if err != nil {
		conf.Logger.Errorf("AutoRun", "error on [newDefaultRipper(conf)]. err:[%s]", err.Error())
		return
	}
	defer r.Close()

	Run(r)
}

// Run runs text processing
func Run(r *Ripper) {
	conf := r.Config
	logger := conf.Logger

	go func() {
		interval := conf.ProgressInterval
		tick := time.Tick(time.Duration(conf.ProgressInterval) * time.Second)
		prev := 0
		for {
			select {
			case <-tick:
				cur := r.GetCurrentPosition()
				logger.Infof("progress", "[%s] line: %d, tps: %d\n", time.Now().Format("2006-01-02 15:04:05"), cur, (cur-prev)/interval)
				prev = cur
			}
		}
	}()

	logger.Infof("Run", "read and write lines...")

	err := r.ReadAndWriteLines()
	if err != nil {
		logger.Errorf("Run", "error on r.ReadAndWriteLines() err:[%s]", err.Error())
		return
	}

	logger.Infof("Run", "finish process")
}

func newDefaultRipper(c Config) (*Ripper, error) {
	r, err := New(c)
	if err != nil {
		return nil, err
	}

	if err := r.WriteHeader(); err != nil {
		r.Close()
		return nil, err
	}

	return r, nil
}
