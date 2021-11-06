package tabulon

import (
	"strings"
	"log"
//	"fmt"
)

type CSVReader struct {
	header []string
	ncols int
	delimiter rune
	quote rune
	columns []string
	column_map []int
	limit int
}

func NewCSVReader() (CSVReader) {
	r := CSVReader{nil, 0, ',', '"', nil, nil, 0}
	return r
}

func (c *CSVReader) SetDelimiter(d rune) {
	c.delimiter = d
}

func (c *CSVReader) SetColumns(cols []string) {
	c.columns = cols
}

func (c *CSVReader) SetLimit(limit int) {
	c.limit = limit
}

func (c *CSVReader) GetHeader() ([]string) {
	return c.header
}

func (c *CSVReader) Reset() {
	c.header = nil
}

func (r *CSVReader) appendToken(row []string, token string) ([]string) {
	if r.limit>0 && len(token)>r.limit {
		token = token[0:r.limit]
	}

	return append(row, token)
}

func (r *CSVReader) tokenize(line string) ([]string) {
	var row []string
	var token string
	i := 0
	N := len(line)

	for {
		// check and handle end of line
		if i>= N {
			break
		}

		// easy case; check if cell is not quoted
		if line[i]!=byte(r.quote) {
			if idx_end := strings.Index(line[i:], string(r.delimiter)); idx_end==-1 {
				token = line[i:N]
				row = r.appendToken(row, token)
				i = N
			} else {
				token = line[i:i+idx_end]
				row = r.appendToken(row, token)
				i += idx_end + 1
			}

		} else {
			i++

			// quoted cell, intermediate position
			quoted_delim := string(r.quote) + string(r.delimiter)
			if idx_end := strings.Index(line[i:], quoted_delim); idx_end!=-1 {
				token = line[i:i+idx_end]
				row = r.appendToken(row, token)
				i += idx_end + 2

				// quoted cell, final position
			} else {
				if idx_end2 := strings.Index(line[i:], string(r.quote)); idx_end2!=-1 {
					token = line[i:i+idx_end2]
					row = r.appendToken(row, token)
					i += idx_end2 + 1

					// quoted cell, unterminated quote
				} else {
					token = line[i:N]
					row = r.appendToken(row, token)
					i = N
				}
			}
		}
	}

	return row
}

func (r *CSVReader) normalizeRow(row []string) ([]string) {
	out := []string{}
	N := len(row)
	//fmt.Printf("norm: input=%v ncols=%v\n", row, r.ncols)
	for i:=0; i<r.ncols; i++ {
		idx := r.column_map[i]
		if idx>=N {
			out = append(out, "")
		} else {
			out = append(out, row[idx])
		}
	}
	return out
}

func findToken(token string, l []string) (int) {
	for i,s := range(l) {
		if token==s {
			return i
		}
	}
	return -1
}

func (r *CSVReader) initializeHeader(row []string) {
	if r.columns == nil {
		r.columns = row
	}

	r.column_map = nil
	for _,col := range(r.columns) {
		idx := findToken(col, row)
		if idx==-1 {
			log.Fatal("specified column not found: ", col)
		}
		r.header = append(r.header, col)
		r.column_map = append(r.column_map, idx)
	}

	if len(r.header)==0 {
		log.Fatal("None of the specified columns were found in the data")
	}

	r.ncols = len(r.header)
}

func (r *CSVReader) ParseLine(line string) (row []string) {
	row = r.tokenize(line)

	if len(r.header) == 0 {
		r.initializeHeader(row)
		return nil
	}

	row = r.normalizeRow(row)
	return row
}
