package main

import (
	"github.com/avraam311/gogrep/cmd/app"
	"github.com/avraam311/gogrep/internal/cutter"
	"github.com/avraam311/gogrep/internal/flags"
)

func main() {
	flags := flags.New()
	cutter := cutter.New()
	app := app.New(cutter, flags)
	app.Run()
}
