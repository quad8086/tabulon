package tabulon

import (
	"os"
	"strings"
	"log"
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

func (table* Table) RenderCSV(fname string) {
	fd,err := os.Create(fname)
	if err!=nil {
		log.Fatal(err)
	}
	
	fd.Write([]byte(strings.Join(table.header, ",") + "\n"))
	for _,row := range(table.content) {
		fd.Write([]byte(strings.Join(row, ",") + "\n"))
	}
	fd.Close()
}
