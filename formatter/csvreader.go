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
}

func NewCSVReader() (CSVReader) {
	r := CSVReader{nil, 0, ',', '"', nil, nil}
	return r
}

func (c *CSVReader) SetDelimiter(d rune) {
	c.delimiter = d
}

func (c *CSVReader) SetColumns(cols []string) {
	c.columns = cols
}

func (c *CSVReader) GetHeader() ([]string) {
	return c.header
}

func (c *CSVReader) Reset() {
	c.header = nil
}

func tokenize(line string, delimiter rune, quote rune) ([]string) {
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
		if line[i]!=byte(quote) {
			if idx_end := strings.Index(line[i:], string(delimiter)); idx_end==-1 {
				token = line[i:N]
				row = append(row, token)
				i = N
			} else {
				token = line[i:i+idx_end]
				row = append(row, token)
				i += idx_end + 1
			}
			
		} else {
			i++
			
			// quoted cell, intermediate position
			quoted_delim := string(quote) + string(delimiter)
			if idx_end := strings.Index(line[i:], quoted_delim); idx_end!=-1 {
				token = line[i:i+idx_end]
				row = append(row, token)
				i += idx_end + 2

				// quoted cell, final position
			} else {
				if idx_end2 := strings.Index(line[i:], string(quote)); idx_end2!=-1 {
					token = line[i:i+idx_end2]
					row = append(row, token)
					i += idx_end2 + 1

					// quoted cell, unterminated quote
				} else {
					token = line[i:N]
					row = append(row, token)
					i = N
				}
			}
		}
	}

	return row
}

func (reader *CSVReader) normalizeRow(row []string) ([]string) {
	out := []string{}
	N := len(row)
	//fmt.Printf("norm: input=%v ncols=%v\n", row, reader.ncols)
	for i:=0; i<reader.ncols; i++ {
		idx := reader.column_map[i]
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

func (reader *CSVReader) initializeHeader(row []string) {
	if reader.columns == nil {
		reader.columns = row
	}

	reader.column_map = nil
	for _,col := range(reader.columns) {
		idx := findToken(col, row)
		if idx==-1 {
			log.Fatal("specified column not found: ", col)
		}
		reader.header = append(reader.header, col)
		reader.column_map = append(reader.column_map, idx)
	}

	if len(reader.header)==0 {
		log.Fatal("None of the specified columns were found in the data")
	}

	reader.ncols = len(reader.header)
}

func (reader *CSVReader) ParseLine(line string) (row []string) {
	row = tokenize(line, reader.delimiter, reader.quote)

	if len(reader.header) == 0 {
		reader.initializeHeader(row)
		return nil
	}

	row = reader.normalizeRow(row)
	return row
}
