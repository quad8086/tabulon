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
		Stdin bool `short:"S" long:"stdin" description:"read from stdin rather than files"`
		Limit int `short:"L" long:"limit" description:"limit cell content length to N"`
		Match []string `short:"m" long:"match" description:"match string (AND)"`
		Expr string `short:"e" long:"expr" description:"match on expression"`
		Plain bool `short:"p" long:"plain" description:"render to stdout as plaintext"`
		CSV bool `short:"C" long:"csv" description:"render to stdout as csv"`
		Skip int `short:"s" long:"skip" description:"skip N lines before load" default:"0"`
		Delimiter string `short:"d" long:"delimiter" description:"set input delimiter" default:""`
		OutputDelimiter string `short:"D" long:"output-delimiter" description:"set output delimiter" default:""`
		Head int `short:"h" long:"head" description:"only consume N first lines of input" default:"-1"`
		Tail int `short:"t" long:"tail" description:"only consume N last lines of input" default:"-1"`
		List string `short:"l" long:"list-column" description:"output specified column as list" default:""`
		Unique string `short:"u" long:"unique" description:"output unique values of specified column as list" default:""`
		TSV bool `long:"tsv" description:"force input delimiter to tab"`
		PSV bool `long:"psv" description:"force input delimiter to pipe"`
		SortColumn string `long:"sort-column" description:"sort by column" default:""`
		Reverse bool `long:"reverse" description:"reverse sort direction"`
		Columns []string `short:"c" long:"columns" description:"only render specified columns"`
	}

	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		return
	}

	files:=args[1:]
	if len(files)==0 && !opts.Stdin {
		log.Fatal("no input provided; please supply filenames or enable stdin")
	}

	if opts.Head>0 && opts.Tail>0 {
		log.Fatal("both head and tail specified; they are mutually exclusive")
	}

	table := tabulon.NewTable()
	table.SetMatch(opts.Match)
	table.SetSkip(opts.Skip)
	table.SetHead(opts.Head)
	table.SetTail(opts.Tail)
	table.SetColumns(opts.Columns)
	table.SetLimit(opts.Limit)
	
	if len(opts.Expr)>0 {
		table.SetMatchExpr(opts.Expr)
	}

	if len(opts.Delimiter)>0 {
		table.SetDelimiter(rune(opts.Delimiter[0]))
	} else if opts.TSV {
		table.SetDelimiter('\t')
	} else if opts.PSV {
		table.SetDelimiter('|')
	}

	if(len(opts.OutputDelimiter)>0) {
		table.SetOutputDelimiter(rune(opts.OutputDelimiter[0]))
	}

	if opts.Stdin {
		table.ReadStdin()
	} else {
		table.ReadFiles(files)
	}

	if len(opts.SortColumn)>0 {
		idx := table.FindColumn(opts.SortColumn)
		if idx==-1 {
			fmt.Println("No such column="+opts.SortColumn)
			os.Exit(1)
		}

		if opts.Reverse {
			table.SortByIndexReverse(idx)
		} else {
			table.SortByIndex(idx)
		}
	}

	if len(opts.List)>0 {
		table.RenderList(opts.List)
	} else if len(opts.Unique)>0 {
		table.RenderUnique(opts.Unique)
	} else if opts.CSV {
		table.RenderCSV()
	} else if opts.Plain || opts.Stdin {
		table.RenderPlaintext()
	} else {
		table.RenderInteractive()
	}
	os.Exit(0)
}
