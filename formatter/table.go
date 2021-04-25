package tabulon

import (
	"os"
	"log"
	"strings"
	"fmt"
	"io"
	"encoding/csv"
	"path"
)

type Table struct {
	header []string
	content [][]string
	match []string
	delimiter rune
	nrows int
	ncols int
	limits []int
	description string
	skip int
}

func NewTable() (Table) {
	t := Table{
		delimiter: ',',
		header: nil,
		content: nil,
		nrows: 0,
		ncols: 0,
	}

	return t
}

func (table* Table) SetMatch(m []string) {
	table.match = m
}

func (table* Table) SetSkip(skip int) {
	table.skip = skip
}

func (table* Table) Clear() {
	table.header = nil
	table.content = nil
	table.nrows = 0
	table.ncols = 0
	table.skip = 0
}

func filter_record(rec []string, match []string) (bool) {
	if len(match) == 0 {
		return false
	}

	line := strings.Join(rec, ",")
	for _,m := range(match) {
		if !strings.Contains(line, m) {
			return true
		}
	}

	return false
}

func (table* Table) processFile(r* csv.Reader) {
	skip := table.skip
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if skip>0 {
			skip--
			continue
		}

		if table.header==nil {
			table.header = rec
			table.ncols = len(table.header)
			continue
		}

		if len(rec) != table.ncols {
			fmt.Printf("skipping record %s\n", rec)
			for _,a := range(rec) {
				fmt.Println(a)
			}
			continue
		}

		if filter_record(rec, table.match) {
			continue
		}

		table.content = append(table.content, rec)
	}
	table.nrows = len(table.content)
}

func (table* Table) ReadStdin() {
	table.Clear()
	table.description = "stdin"
	r := csv.NewReader(os.Stdin)
	table.processFile(r)
	table.calcLimits()
}

func (table* Table) ReadFiles(files []string) {
	if len(files)==0 {
		log.Fatal("ReadFiles: no files to read")
	}

	table.header = nil
	table.description = path.Base(files[0])

	for _,file := range files {
		fd, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		r := csv.NewReader(fd)
		r.Comma = table.delimiter
		table.processFile(r)
	}
	table.calcLimits()
}

func (table* Table) calcLimits() {
	ncols := len(table.header)
	table.limits = make([]int, ncols)
	for j,cell := range(table.header) {
		table.limits[j] = int_max(table.limits[j], len(cell))
	}

	for _,row := range(table.content) {
		for j,cell := range(row) {
			table.limits[j] = int_max(table.limits[j], len(cell))
		}
	}
}
