package tabulon

import (
	"os"
	"strings"
	"log"
	"fmt"
)

func (table *Table) RenderPlaintext() {
	header_print(table.header, table.limits)
	for _,row := range(table.content) {
		row_print(row, table.limits)
	}
}

func (table *Table) RenderInteractive() {
	term := NewTerminal()
	term.Run(table)
}

func (table *Table) RenderCSV() {
	fd := os.Stdout
	d := string(table.output_delimiter)
	fd.Write([]byte(strings.Join(table.header, d) + "\n"))
	for _,row := range(table.content) {
		fd.Write([]byte(strings.Join(row, d) + "\n"))
	}
	fd.Close()
}

func (table *Table) RenderList(col string) {
	idx := table.FindColumn(col)
	if idx == -1 {
		log.Fatal("RenderList: no such column="+col)
	}

	var output string
	for _,row := range(table.content) {
		if len(output)>0 {
			output += string(table.output_delimiter)
		}

		output += row[idx]
	}

	fmt.Println(output)
}

func (table *Table) RenderUnique(col string) {
	idx := table.FindColumn(col)
	if idx == -1 {
		log.Fatal("RenderUnique: no such column="+col)
	}

	set := make(map[string]bool)
	for _,row := range(table.content) {
		set[row[idx]] = true
	}
	var output string
	for k,_ := range(set) {
		if len(output)>0 {
			output += string(table.output_delimiter)
		}

		output += k
	}

	fmt.Println(output)
}
