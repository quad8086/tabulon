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
		Delimiter string `short:"d" long:"delimiter" description:"delimiter"`
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

	var table Table
	table.match = opts.Match
	table.delimiter = opts.Delimiter
	
	if opts.Stdin {
		table.ReadStdin()
	} else {
		table.ReadFiles(files)
	}

	if opts.Plain {
		table.RenderPlain()
	} else {
		table.RenderTCell()
	}
	os.Exit(0)
}
