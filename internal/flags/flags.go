package flags

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type Flags struct {
	FlagF    string
	FlagD    string
	FlagS    bool
	Filename string
}

func New() *Flags {
	FlagF := pflag.StringP("f", "f", "", "columns to print")
	FlagD := pflag.StringP("d", "d", "\t", "change delimiter")
	FlagS := pflag.BoolP("s", "s", false, "ignore strings withour delimiter")
	pflag.Parse()
	fileName := pflag.Arg(0)

	if *FlagF == "" {
		fmt.Println("flag -f is required")
		os.Exit(1)
	}

	flags := Flags{
		FlagF:    *FlagF,
		FlagD:    *FlagD,
		FlagS:    *FlagS,
		Filename: fileName,
	}

	return &flags
}
