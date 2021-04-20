package tabulon

import (
	"log"
	"os"
	"fmt"
	"github.com/gdamore/tcell/v2"
)

type Terminal struct {
	have_header bool
	have_status bool
	yheader int
	ystatus int
	xscreen int
	yscreen int
	screen tcell.Screen
}

func NewTerminal() (Terminal) {
	t := Terminal{
		have_header: true,
		have_status: true,
		yheader: 0,
		ystatus: 0,
	}

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

	t.xscreen, t.yscreen = s.Size()
	t.screen = s
	return t
}

func tcell_line(s tcell.Screen, x, y int, line string, style tcell.Style) {
	for i, c := range(line) {
		s.SetContent(x+i, y, c, nil, style)
	}
}

func tcell_row(s tcell.Screen, x, y int, row []string, lim []int, style tcell.Style) {
	for i,cell := range(row) {
		for k,r := range(cell) {
			xpos := x+k
			s.SetContent(xpos, y, r, nil, style)
		}
		x += (lim[i] + 1)
	}
}

func tcell_render(table* Table, term* Terminal, yview int) (int) {
	term.screen.Clear()

	ylim_hi := table.nrows - term.yscreen
	yview = int_max(0, yview)	
	if yview>ylim_hi {
		yview = ylim_hi
	}

	style_normal := tcell.StyleDefault
	style_underl := tcell.StyleDefault.Underline(true)
	y_header := 0
	y_status := term.yscreen-1
	tcell_row(term.screen, 0, y_header, table.header, table.limits, style_underl)

	content_index := 0
	for y:=1; y<y_status; y++ {
		tcell_row(term.screen, 0, y, table.content[yview + content_index], table.limits,
			style_normal)
		content_index++
	}
	
	tcell_line(term.screen, 0, y_status, fmt.Sprintf("%s: nrows=%d ncols=%d yview=%d",
		table.description, table.nrows, table.ncols, yview), style_underl)
	term.screen.Show()
	return yview
}

func (term* Terminal) Run(table* Table) {	
	yview := 0
	s:= term.screen
	
	for {
		yview = tcell_render(table, term, yview)
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			_, term.yscreen = s.Size()
			
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Rune()=='q' {
				s.Fini()
				os.Exit(0)
			}

			if ev.Key() == tcell.KeyDown || ev.Rune()=='j' {
				yview = yview+1
			}

			if ev.Key() == tcell.KeyUp || ev.Rune()=='k' {
				yview = yview-1
			}

			if ev.Key() == tcell.KeyPgDn || ev.Rune()==' ' {
				yview = yview+term.yscreen
			}

			if ev.Key() == tcell.KeyPgUp || ev.Rune()=='b' {
				yview = yview-term.yscreen
			}

			if ev.Key() == tcell.KeyHome || ev.Rune()=='0' {
				yview = 0
			}

			if ev.Key() == tcell.KeyEnd || ev.Rune()=='G' {
				yview = table.nrows - term.yscreen
			}

			if ev.Rune() == 'S' {
				term.have_status = !term.have_status
			}

			if ev.Rune() == 'H' {
				term.have_header = !term.have_header
			}
		}
	}
}
