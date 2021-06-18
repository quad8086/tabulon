package tabulon

import (
	"os"
	"log"
	"strings"
	"bufio"
	"path"
	"sort"
)

type Table struct {
	header []string
	content [][]string
	match []string
	delimiter rune
	output_delimiter rune
	nrows int
	ncols int
	limits []int
	description string
	skip int
	head int
	tail int
}

func NewTable() (Table) {
	t := Table {
		delimiter: 0,
		header: nil,
		content: nil,
		nrows: 0,
		ncols: 0,
		output_delimiter: ',',
		skip: 0,
		head: -1,
		tail: -1,
	}

	return t
}

func (table *Table) Clear() {
	table.header = nil
	table.content = nil
	table.nrows = 0
	table.ncols = 0
	table.skip = 0
	table.head = -1
	table.tail = -1
}

func (table *Table) SetMatch(m []string) {
	table.match = m
}

func (table *Table) SetSkip(skip int) {
	table.skip = skip
}

func (table *Table) SetHead(head int) {
	if head< -1 {
		log.Fatal("head has to be positive")
	}
	table.head = head
}

func (table *Table) SetTail(tail int) {
	if tail< -1 {
		log.Fatal("tail has to be positive")
	}
	table.tail = tail
}

func (table *Table) SetDelimiter(d rune) {
	table.delimiter = d
}

func (table *Table) SetOutputDelimiter(d rune) {
	table.output_delimiter = d
}

func (table *Table) calcLimits() {
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

func filter_record(rec []string, t *Table) (bool) {
	if len(t.match) == 0 {
		return false
	}

	line := strings.Join(rec, string(t.delimiter))
	for _,m := range(t.match) {
		if !strings.Contains(line, m) {
			return true
		}
	}

	return false
}

func (table *Table) processFile(fd *os.File) {
	skip := table.skip
	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanLines)
	reader := NewCSVReader()
	reader.SetDelimiter(table.delimiter)
	head := table.head
	tail := table.tail

	for scanner.Scan() {
		if skip>0 {
			skip--
			continue
		}

		row := reader.ParseLine(scanner.Text())
		if table.header==nil {
			table.header = reader.GetHeader()
			table.ncols = len(table.header)
			continue
		}

		if filter_record(row, table) {
			continue
		}

		if head!=-1 {
			if head==0 {
				break
			}

			head--
		}

		table.content = append(table.content, row)
	}

	n := len(table.content)
	if tail!=-1 && n>tail {
		table.content = table.content[(n-tail):n]
	}

	table.nrows = len(table.content)
}

func (table *Table) ReadStdin() {
	table.Clear()
	table.description = "stdin"
	if table.delimiter==0 {
		table.delimiter = ','
	}
	table.processFile(os.Stdin)
	table.calcLimits()
}

func guess_delimiter(fname string) (rune) {
	if strings.Contains(fname, ".csv") {
		return ','
	} else if strings.Contains(fname, ".psv") {
		return '|'
	} else if strings.Contains(fname, ".tsv") {
		return '\t'
	}

	return ','
}

func (table *Table) ReadFiles(files []string) {
	if files==nil || len(files)==0 {
		log.Fatal("ReadFiles: no files to read")
	}

	table.header = nil
	table.description = path.Base(files[0])

	for _,file := range files {
		fd, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		if table.delimiter==0 {
			table.delimiter = guess_delimiter(file)
		}

		table.processFile(fd)
		fd.Close()
	}
	table.calcLimits()
}

func (table *Table) Search(yorig int, s string) (int) {
	for y:=yorig+1; y<len(table.content); y++ {
		row := table.content[y]
		for _,cell := range(row) {
			if(strings.Contains(cell, s)) {
				return y
			}
		}
	}

	return yorig
}

func (table *Table) SearchReverse(yorig int, s string) (int) {
	for y:=yorig-1; y>=0; y-- {
		row := table.content[y]
		for _,cell := range(row) {
			if(strings.Contains(cell, s)) {
				return y
			}
		}
	}

	return yorig
}

func (table *Table) FindColumn(col string) (int) {
	for i, s := range(table.header) {
		if s==col {
			return i
		}
	}

	return -1
}

func (table *Table) SortByIndex(idx int) {
	sort.Slice(table.content, func(i, j int) bool {
		return table.content[i][idx] < table.content[j][idx]
	})
}

func (table *Table) SortByIndexReverse(idx int) {
	sort.Slice(table.content, func(i, j int) bool {
		return table.content[j][idx] < table.content[i][idx]
	})
}
