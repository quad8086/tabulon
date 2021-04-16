package tabulon

import (
	"fmt"
	"log"
	"os"
	"github.com/gdamore/tcell/v2"
)

func int_max(i int, j int) (int) {
	if(i>j) {
		return i
	}
	return j
}

func int_min(i int, j int) (int) {
	if(i>j) {
		return j
	}
	return i
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

func gen_padded_cell(cell string, lim int) (string) {
	res := cell
	for i:=len(cell); i<lim; i++ {
		res = res + " ";
	}

	res += " |"
	return res
}

func row_print(row []string, limits []int) {
	var line string
	for i,cell := range(row) {
		line += gen_padded_cell(cell, limits[i])
	}

	fmt.Print(line+"\n")
}

func header_print(row []string, limits []int) {
	row_print(row, limits)

	var line string
	for i,_ := range(row) {
		for l:=0; l<limits[i]; l++ {
			line += "-";
		}
		line += "-+";
	}
	fmt.Print(line + "\n")
}

func (table* Table) RenderPlain() {
	table.calcLimits()
	header_print(table.header, table.limits)
	for _,row := range(table.content) {
		row_print(row, table.limits)
	}
}

func tcell_line(s tcell.Screen, x, y int, line string) {
	for i, c := range(line) {
		s.SetContent(x+i, y, c, nil, tcell.StyleDefault.Underline(true))
	}
}

func tcell_render(s tcell.Screen, table* Table, yoffset int) (int,int) {
	screenx, screeny := s.Size()
	s.Clear()

	ylim_lo := 0
	ylim_hi := table.nrows - screeny
	
	if yoffset<0 {
		yoffset = ylim_lo
	}

	if yoffset>ylim_hi {
		yoffset = ylim_hi
	}

	for y:=0; y<screeny-1; y++ {
		row := table.content[yoffset+y]
		x := 0
		for i,cell := range(row) {
			for k,r := range(cell) {
				xpos := x+k
				if xpos>screenx-1 {
					break
				}
				
				s.SetContent(xpos, y, r, nil, tcell.StyleDefault)
			}
			x += table.limits[i]
		}
	}

	tcell_line(s, 0, screeny-1, fmt.Sprintf("%d/%d", yoffset, ylim_hi))
	s.Show()

	return yoffset, screeny
}

func (table* Table) RenderTCell() {
	s, e := tcell.NewScreen()
	if e != nil {
		log.Fatal(e)
	}
	if e := s.Init(); e != nil {
		log.Fatal(e)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	_, screeny := s.Size()
	yoffset := 0
	_,yscreen := s.Size()
	
	for {
		yoffset, yscreen = tcell_render(s, table, yoffset)
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			_, screeny = s.Size()
			
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Rune()=='q' {
				s.Fini()
				os.Exit(0)
			}

			if ev.Key() == tcell.KeyDown {
				yoffset = yoffset+1
			}

			if ev.Key() == tcell.KeyUp {
				yoffset = yoffset-1
			}

			if ev.Key() == tcell.KeyPgDn {
				yoffset = yoffset+screeny
			}

			if ev.Key() == tcell.KeyPgUp {
				yoffset = yoffset-screeny
			}

			if ev.Key() == tcell.KeyHome {
				yoffset = 0
			}

			if ev.Key() == tcell.KeyEnd {
				yoffset = table.nrows - yscreen
			}
		}
	}
}
