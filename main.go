package main

import (
	"log"
	"time"

	"github.com/mmrezoe/tasks/config"
	"github.com/mmrezoe/tasks/runner"
)

func main() {
	c := config.ParseArgs()

	startTime := time.Now()
	config.Debug("start: "+time.Now().Format("2006-01-02 15:04:05"), c.Debug)

	err := runner.Run(c)
	if err != nil {
		log.Fatal(err)
	}

	config.Debug("time elapsed: "+time.Now().Sub(startTime).String(), c.Debug)
}
