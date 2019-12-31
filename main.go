package main

import (
	"fmt"

	"github.com/evalphobia/go-jp-text-ripper/prefilter"
	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

var (
	version  string
	revision string
)

// cli entry point
func main() {
	conf := ripper.Config{
		Version:  version,
		Revision: revision,
	}
	if err := conf.Init(); err != nil {
		fmt.Printf("error on conf.Init(). err:[%s]", err.Error())
		return
	}
	if conf.UseNeologd {
		conf.PreFilters = append(conf.PreFilters, prefilter.Neologd)
	}

	ripper.AutoRun(conf)
}
