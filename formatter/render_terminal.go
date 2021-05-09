package tabulon

import (
	"log"
	"os"
	"fmt"
	"strconv"
	"github.com/gdamore/tcell/v2"
)

type UIMode int
const (
	Normal UIMode = iota
	Search
	SearchReverse
)

type Terminal struct {
	yheader int
	ystatus int
	xscreen int
	yscreen int
	screen tcell.Screen
	mode UIMode
	search string
	style_normal tcell.Style
	style_underl tcell.Style
}

func NewTerminal() (Terminal) {
	t := Terminal{
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
	t.style_normal = tcell.StyleDefault
	t.style_underl = tcell.StyleDefault.Underline(true)
	return t
}

func tcell_line(s tcell.Screen, x, y int, line string, style tcell.Style) {
	for i, c := range(line) {
		s.SetContent(x+i, y, c, nil, style)
	}
}

func tcell_row(s tcell.Screen, x, y int, xstart int, row []string, lim []int, style tcell.Style) {
	for i,cell := range(row) {
		if i<xstart {
			continue
		}
		for k,r := range(cell) {
			xpos := x+k
			s.SetContent(xpos, y, r, nil, style)
		}
		x += (lim[i] + 1)
	}
}

func tcell_render(table* Table, term* Terminal, xview, yview int) (int, int) {
	term.screen.Clear()

	xlim_hi := table.ncols-1
	ylim_hi := table.nrows-1
	yview = int_max(0, yview)
	xview = int_max(0, xview)
	if yview>ylim_hi {
		yview = ylim_hi
	}
	if xview>xlim_hi {
		xview = xlim_hi
	}
	
	y_header := 0
	y_status := term.yscreen-1
	tcell_row(term.screen, 0, y_header, xview, table.header, table.limits, term.style_underl)

	empty_row := []string{"~"}
	content_index := 0
	for y:=1; y<y_status; y++ {
		idx := yview + content_index
		row := empty_row
		if idx>=0 && idx<=ylim_hi {
			row = table.content[idx]
		}
		
		tcell_row(term.screen, 0, y, xview, row, table.limits, term.style_normal)
		content_index++
	}

	if(term.mode == Normal) {
		tcell_line(term.screen, 0, y_status,
			fmt.Sprintf("%v: row=%v/%v col=%v/%v screen=%v,%v",
				table.description,
				yview, table.nrows,
				xview, table.ncols,
				term.yscreen, term.xscreen),
			term.style_underl)
		
	} else if(term.mode == Search) {
		tcell_line(term.screen, 0, y_status, fmt.Sprintf("Search: %v", term.search),
			term.style_underl)
		
	} else if(term.mode == SearchReverse) {
		tcell_line(term.screen, 0, y_status, fmt.Sprintf("Search reverse: %v", term.search),
			term.style_underl)
	}
	
	term.screen.Show()
	return xview, yview
}

func (term* Terminal) Run(table* Table) {	
	yview := 0
	xview := 0
	s:= term.screen
	
	for {
		xview, yview = tcell_render(table, term, xview, yview)
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			_, term.yscreen = s.Size()
			
		case *tcell.EventKey:
			if(term.mode == Normal) {
				if ev.Key() == tcell.KeyEscape || ev.Rune()=='q' || ev.Rune()=='Q' {
					s.Fini()
					os.Exit(0)
				}
				
				if ev.Key() == tcell.KeyDown || ev.Rune()=='j' {
					yview = yview+1
				}
				
				if ev.Key() == tcell.KeyUp || ev.Rune()=='k' {
					yview = yview-1
				}

				if ev.Key() == tcell.KeyLeft || ev.Rune()=='h' {
					xview = xview-1
				}

				if ev.Key() == tcell.KeyRight || ev.Rune()=='l' {
					xview = xview+1
				}
				
				if ev.Key() == tcell.KeyPgDn || ev.Rune()==' ' {
					yview = yview+term.yscreen
				}
				
				if ev.Key() == tcell.KeyPgUp || ev.Rune()=='b' {
					yview = yview-term.yscreen
				}
				
				if ev.Key() == tcell.KeyHome || ev.Rune()=='0' || ev.Rune()=='g' {
					yview = 0
					xview = 0
				}
				
				if ev.Key() == tcell.KeyEnd || ev.Rune()=='G' {
					yview = table.nrows - term.yscreen
					xview = 0
				}
				
				if ev.Rune() == '/' {
					term.mode = Search
				}

				if ev.Rune() == '?' {
					term.mode = SearchReverse
				}
				
				if ev.Rune() == 'n' && len(term.search)>0 {
					yview = table.Search(yview, term.search)
				}

				if ev.Rune() == 'N' && len(term.search)>0 {
					yview = table.SearchReverse(yview, term.search)
				}
				
			} else if(term.mode == Search || term.mode==SearchReverse) {
				if term.mode==Search && ev.Key() == tcell.KeyEnter {
					yview = table.Search(yview, term.search)
					term.mode = Normal
				}

				if term.mode==SearchReverse && ev.Key() == tcell.KeyEnter {
					yview = table.SearchReverse(yview, term.search)
					term.mode = Normal
				}

				if ev.Key() == tcell.KeyCtrlG {
					term.search = ""
					term.mode = Normal
				}

				if ev.Key() == tcell.KeyCtrlU {
					term.search = ""
				}

				is_backspace := ev.Key() == tcell.KeyBackspace ||
					ev.Key() == tcell.KeyDelete ||
					ev.Key() == tcell.KeyCtrlH ||
					ev.Key() == tcell.KeyBackspace2
				if is_backspace && len(term.search)>0 {
					term.search = term.search[:len(term.search)-1]
				}

				if strconv.IsPrint(ev.Rune()) {
					term.search += string(ev.Rune())
				}
			}
		}
	}
}
