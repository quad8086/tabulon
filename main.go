package main

import (
	"os"
	"log"
	"github.com/jessevdk/go-flags"
	"tabulon/formatter"
	"fmt"
)

func main() {
	var opts struct {
		Stdin bool `short:"S" long:"stdin" description:"read from stdin"`
		Match []string `short:"m" long:"match" description:"match string"`
		Plain bool `short:"p" long:"plain" description:"render to stdout as plaintext"`
		CSV bool `short:"C" long:"csv" description:"render as csv file to stdout"`
		Skip int `short:"s" long:"skip" description:"skip N lines before load" default:"0"`
		Delimiter string `short:"d" long:"delimiter" description:"set delimiter" default:""`
		OutputDelimiter string `short:"D" long:"output-delimiter" description:"set output delimiter" default:""`
	}

	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	files:=args[1:]
	if len(files)==0 && !opts.Stdin {
		log.Fatal("no input provided; please supply filenames or enable stdin")
	}

	table := tabulon.NewTable()
	table.SetMatch(opts.Match)
	table.SetSkip(opts.Skip)

	if(len(opts.Delimiter)>0) {
		table.SetDelimiter(rune(opts.Delimiter[0]))
	}

	if(len(opts.OutputDelimiter)>0) {
		table.SetOutputDelimiter(rune(opts.OutputDelimiter[0]))
	}

	if opts.Stdin {
		table.ReadStdin()
	} else {
		table.ReadFiles(files)
	}

	if opts.CSV {
		table.RenderCSV()
	} else if opts.Plain || opts.Stdin {
		table.RenderPlaintext()
	} else {
		table.RenderInteractive()
	}
	os.Exit(0)
}
