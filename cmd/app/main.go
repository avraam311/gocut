package app

import (
	"github.com/avraam311/gogrep/internal/cutter"
	"github.com/avraam311/gogrep/internal/flags"
)

type App struct {
	cutter *cutter.Cutter
	flags  *flags.Flags
}

func New(cut *cutter.Cutter, flags *flags.Flags) *App {
	return &App{
		cutter: cut,
		flags:  flags,
	}
}

func (a *App) Run() {
	a.cutter.Cut(a.flags.FlagD, a.flags.FlagF, a.flags.Filename, a.flags.FlagS)
}
