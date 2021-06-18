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
	Help
)

type Terminal struct {
	xscreen int
	yscreen int
	screen tcell.Screen
	mode UIMode
	search string
	style_normal tcell.Style
	style_underl tcell.Style
	xview int
	yview int
	help []string
}

func NewTerminal() (Terminal) {
	t := Terminal {
		mode: Normal,
		xview: 0,
		yview: 0,
	}

	s, e := tcell.NewScreen()
	if e != nil {
		log.Fatal(e)
	}
	if e := s.Init(); e != nil {
		log.Fatal(e)
	}

	t.xscreen, t.yscreen = s.Size()
	t.screen = s
	t.style_normal = tcell.StyleDefault
	t.style_underl = t.style_normal.Underline(true)

	t.help = []string{"help1", "help2"}
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

func tcell_render(term *Terminal, table *Table) {
	term.screen.Clear()

	xlim_hi := table.ncols-1
	ylim_hi := table.nrows-1
	term.yview = int_max(0, term.yview)
	term.xview = int_max(0, term.xview)
	if term.yview>ylim_hi {
		term.yview = ylim_hi
	}
	if term.xview>xlim_hi {
		term.xview = xlim_hi
	}

	y_header := 0
	y_status := term.yscreen-1

	if(term.mode == Help) {
		tcell_row(term.screen, 0, y_header, term.xview, []string{"Help"}, []int{5}, term.style_underl)
	} else {
		tcell_row(term.screen, 0, y_header, term.xview, table.header, table.limits, term.style_underl)
	}

	empty_row := []string{"~"}
	content_index := 0
	for y:=1; y<y_status; y++ {
		idx := term.yview + content_index
		row := empty_row
		if idx>=0 && idx<=ylim_hi {
			row = table.content[idx]
		}

		tcell_row(term.screen, 0, y, term.xview, row, table.limits, term.style_normal)
		content_index++
	}

	if(term.mode == Normal) {
		tcell_line(term.screen, 0, y_status,
			fmt.Sprintf("%v: row=%v/%v col=%v/%v screen=%v,%v",
				table.description,
				term.yview, table.nrows,
				term.xview, table.ncols,
				term.yscreen, term.xscreen),
			term.style_underl)

	} else if(term.mode == Search) {
		tcell_line(term.screen, 0, y_status, fmt.Sprintf("Search: %v", term.search),
			term.style_underl)

	} else if(term.mode == SearchReverse) {
		tcell_line(term.screen, 0, y_status, fmt.Sprintf("Search reverse: %v", term.search),
			term.style_underl)
	} else if(term.mode == Help) {
		tcell_line(term.screen, 0, y_status, "Help", term.style_underl)
	}

	term.screen.Show()
}

func run_normal(ev *tcell.EventKey, term *Terminal, table *Table) {
	s := term.screen

	if ev.Key() == tcell.KeyEscape || ev.Rune()=='q' || ev.Rune()=='Q' {
		s.Fini()
		os.Exit(0)
	}

	if ev.Key() == tcell.KeyDown || ev.Rune()=='j' {
		term.yview = term.yview+1
	}

	if ev.Key() == tcell.KeyUp || ev.Rune()=='k' {
		term.yview = term.yview-1
	}

	if ev.Key() == tcell.KeyLeft || ev.Rune()=='h' {
		term.xview = term.xview-1
	}

	if ev.Key() == tcell.KeyRight || ev.Rune()=='l' {
		term.xview = term.xview+1
	}

	if ev.Key() == tcell.KeyPgDn || ev.Rune()==' ' {
		term.yview = term.yview+term.yscreen
	}

	if ev.Key() == tcell.KeyPgUp || ev.Rune()=='b' {
		term.yview = term.yview-term.yscreen
	}

	if ev.Rune()=='0' {
		term.yview = 0
		term.xview = 0
	}

	if ev.Key() == tcell.KeyHome || ev.Rune()=='g' {
		term.yview = 0
	}

	if ev.Key() == tcell.KeyEnd || ev.Rune()=='G' {
		term.yview = table.nrows - term.yscreen
	}

	if ev.Rune() == '/' {
		term.mode = Search
	}

	if ev.Rune() == '?' {
		term.mode = SearchReverse
	}

	if ev.Rune() == 'h' {
		term.mode = Help
	}

	if ev.Rune() == 's' {
		table.SortByIndex(term.xview)
	}

	if ev.Rune() == 'S' {
		table.SortByIndexReverse(term.xview)
	}

	if ev.Rune() == 'n' && len(term.search)>0 {
		term.yview = table.Search(term.yview, term.search)
	}

	if ev.Rune() == 'N' && len(term.search)>0 {
		term.yview = table.SearchReverse(term.yview, term.search)
	}
}

func run_search(ev *tcell.EventKey, term *Terminal, table *Table) {
	if term.mode==Search && ev.Key() == tcell.KeyEnter {
		term.yview = table.Search(term.yview, term.search)
		term.mode = Normal
	}

	if term.mode==SearchReverse && ev.Key() == tcell.KeyEnter {
		term.yview = table.SearchReverse(term.yview, term.search)
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

func run_help(ev *tcell.EventKey, term *Terminal, table *Table) {
	if ev.Key() == tcell.KeyEscape || ev.Rune()=='q' || ev.Rune()=='Q' {
		term.mode = Normal
	}
}

func (term *Terminal) Run(table *Table) {
	for {
		tcell_render(term, table)
		switch ev := term.screen.PollEvent().(type) {
		case *tcell.EventResize:
			term.screen.Sync()
			term.xscreen, term.yscreen = term.screen.Size()

		case *tcell.EventKey:
			if(term.mode == Normal) {
				run_normal(ev, term, table)
			} else if(term.mode == Search || term.mode==SearchReverse) {
				run_search(ev, term, table)
			} else if(term.mode == Help) {
				run_help(ev, term, table)
			}
		}
	}
}
