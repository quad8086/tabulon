package tabulon

import (
	"strings"
//	"fmt"
)

type CSVReader struct {
	header []string
	row []string
	ncols int
	delimiter rune
}

func NewCSVReader() (CSVReader) {
	r := CSVReader{nil, nil, 0, ','}
	return r
}

func tokenize(line string, delimiter rune) ([]string) {
	if !strings.Contains(line, `"`) {
		return strings.Split(line, string(delimiter))
	}

	var in_quote bool = false
	var res []string
	var token string
	for _,c := range(line) {
		if c=='\'' || c=='"' {
			in_quote = !in_quote
			continue
		}

		if !in_quote && c==delimiter {
			res = append(res, token)
			token = ""
			continue
		}

		token += string(c)
	}
	if len(token)>0 {
		res = append(res, token)
	}

	return res
}

func normalize(row []string, ncols int) ([]string) {
	n := len(row)
	//fmt.Printf("%d %d\n", n, ncols)
	if n<ncols {
		for i:=n; i<ncols; i++ {
			row = append(row, "")
		}
	} else if n>ncols {
		row = row[0:ncols]
	}

	return row
}

func (c* CSVReader) GetHeader() ([]string) {
	return c.header
}

func (c* CSVReader) Reset() {
	c.header = nil
}

func (reader* CSVReader) ParseLine(line string) (row []string) {
	if len(reader.header) == 0 {
		reader.header = tokenize(line, reader.delimiter)
		reader.ncols = len(reader.header)
		return nil
	}

	reader.row = tokenize(line, reader.delimiter)
	reader.row = normalize(reader.row, reader.ncols)
	return reader.row
}
