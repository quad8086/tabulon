package main

import (
	"os"
	"log"
	"github.com/jessevdk/go-flags"
	"tabulon/formatter"
)

func main() {
	var opts struct {
		Stdin bool `short:"S" long:"stdin" description:"read from stdin"`
		Match []string `short:"m" long:"match" description:"match string"`
		Plain bool `short:"p" long:"plain" description:"render to stdout as plaintext"`
		CSV string `short:"C" long:"csv" description:"render as csv file"`
		Skip int `short:"s" long:"skip" description:"skip N lines before load" default:"0"`
	}

	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	files:=args[1:]
	if len(files)==0 && !opts.Stdin {
		log.Fatal("no input provided; please supply filenames or enable stdin")
	}

	table := tabulon.NewTable()
	table.SetMatch(opts.Match)
	table.SetSkip(opts.Skip)

	if opts.Stdin {
		table.ReadStdin()
	} else {
		table.ReadFiles(files)
	}

	if opts.Plain {
		table.RenderPlaintext()
	} else if len(opts.CSV)>0 {
		table.RenderCSV(opts.CSV)
	} else {
		table.RenderTerminal()
	}
	os.Exit(0)
}
