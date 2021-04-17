package tabulon

import (
	"os"
	"log"
	"github.com/jessevdk/go-flags"
)

func Run() {
	var opts struct {
		Stdin bool `short:"S" long:"stdin" description:"read from stdin"`
		Match []string `short:"m" long:"match" description:"match string"`
		Plain bool `short:"p" long:"plain" description:"dump output, no ui"`
//		Delimiter rune `short:"d" long:"delimiter" description:"delimiter"`
		Skip int `long:"skip" description:"skip initial N lines"`
	}

	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	files:=args[1:]
	if len(files)==0 && !opts.Stdin {
		log.Fatal("no input; please supply filenames or enable stdin")
	}

	table := NewTable()
	table.match = opts.Match
	
	if opts.Stdin {
		table.ReadStdin()
	} else {
		table.ReadFiles(files)
	}

	if opts.Plain {
		table.RenderPlaintext()
	} else {
		table.RenderTerminal()
	}
	os.Exit(0)
}
