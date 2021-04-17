package tabulon

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
