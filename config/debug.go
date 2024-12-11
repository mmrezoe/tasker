package config

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Debug(text string, debug bool) {
	if debug {
		fmt.Fprintf(os.Stderr, "[%s] %s\n", color.HiCyanString("tracker"), text)
	}
}
