package tabulon

import (
	"os"
	"log"
	"strings"
	"bufio"
	"path"
	"sort"
	"strconv"
	"github.com/danielgtaylor/mexpr"
)

type Table struct {
	header []string
	content [][]string
	match []string
	remove []string
	delimiter rune
	output_delimiter rune
	nrows int
	ncols int
	limits []int
	description string
	skip int
	head int
	tail int
	limit int
	columns []string
	match_parser mexpr.Interpreter
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
		limit: 0,
		columns: nil,
		match_parser: nil,
	}

	return t
}

func (table *Table) IsEmpty() bool {
	return len(table.content) == 0
}

func (table *Table) SetMatch(m []string) {
	table.match = m
}

func (table *Table) SetRemove(m []string) {
	table.remove = m
}

func (table *Table) SetSkip(skip int) {
	table.skip = skip
}

func (table *Table) SetLimit(limit int) {
	table.limit = limit
}

func (table *Table) SetMatchExpr(match_expr string) {
	l := mexpr.NewLexer(match_expr)
	p := mexpr.NewParser(l)
	ast, err := p.Parse()
	if err != nil {
		log.Fatal("filterRow: invalid expression: ", match_expr)
	}

	table.match_parser = mexpr.NewInterpreter(ast)
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

func (table *Table) SetColumns(c []string) {
	table.columns = c
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

func acceptRow(rec []string, t *Table) (bool) {
	if len(t.remove) > 0 {
		line := strings.Join(rec, string(t.delimiter))
		for _,m := range(t.remove) {
			if strings.Contains(line, m) {
				return false
			}
		}
	}

	if len(t.match) > 0 {
		line := strings.Join(rec, string(t.delimiter))
		for _,m := range(t.match) {
			if !strings.Contains(line, m) {
				return false
			}
		}
	}

	if t.match_parser != nil {
		vars := make(map[string]interface{})
		for i,h := range(t.header) {
			v, err := strconv.ParseFloat(rec[i], 32)
			if err == nil {
				vars[h] = v
			} else {
				vars[h] = rec[i]
			}
		}

		result, err := t.match_parser.Run(vars)
		if err == nil && result==false {
			return false
		}
	}

	return true
}

func (table *Table) processFile(fd *os.File) {
	skip := table.skip
	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanLines)
	reader := NewCSVReader()
	reader.SetDelimiter(table.delimiter)
	reader.SetColumns(table.columns)
	reader.SetLimit(table.limit)
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

		if !acceptRow(row, table) {
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

func guessDelimiter(fname string) (rune) {
	fname = strings.ToLower(fname)
	if strings.Contains(fname, ".csv") {
		return ','
	} else if strings.Contains(fname, ".psv") {
		return '|'
	} else if strings.Contains(fname, ".tsv") {
		return '\t'
	}

	return ','
}

func (table *Table) ReadStdin() {
	table.description = "stdin"
	if table.delimiter==0 {
		table.delimiter = ','
	}

	table.processFile(os.Stdin)
	table.calcLimits()
}

func (table *Table) ReadFiles(files []string) {
	if files==nil || len(files)==0 {
		log.Fatal("ReadFiles: no files to read")
	}

	table.header = nil
	table.description = path.Base(files[0])

	for _,file := range files {
		fi, err := os.Stat(file)
		if err != nil && !fi.Mode().IsRegular() {
			continue
		}

		fd, err := os.Open(file)
		if err != nil {
			continue
		}

		if table.delimiter==0 {
			table.delimiter = guessDelimiter(file)
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

func (table *Table) performSort(idx int, rev bool) {
	sort.Slice(table.content, func(i, j int) bool {
		s1 := table.content[i][idx]
		s2 := table.content[j][idx]
		f1, err1 := strconv.ParseFloat(s1, 8)
		if err1 == nil {
			f2, err2 := strconv.ParseFloat(s2, 8)
			if err2 == nil {
				if rev {
					return f1 > f2
				} else {
					return f1 < f2
				}
			}
		}

		if rev {
			return s1 > s2
		} else {
			return s1 < s2
		}
	})
}

func (table *Table) SortByIndex(idx int) {
	table.performSort(idx, false)
}

func (table *Table) SortByIndexReverse(idx int) {
	table.performSort(idx, true)
}
