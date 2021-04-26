package tabulon

import (
	"os"
	"strings"
)

func (table* Table) RenderPlaintext() {
	header_print(table.header, table.limits)
	for _,row := range(table.content) {
		row_print(row, table.limits)
	}
}

func (table* Table) RenderTerminal() {
	term := NewTerminal()
	term.Run(table)
}

func (table* Table) RenderCSV() {
	fd := os.Stdout
	d := string(table.output_delimiter)
	fd.Write([]byte(strings.Join(table.header, d) + "\n"))
	for _,row := range(table.content) {
		fd.Write([]byte(strings.Join(row, d) + "\n"))
	}
	fd.Close()
}
