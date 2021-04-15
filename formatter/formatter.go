package tdf

import (
	"os"
	"fmt"
	"log"
	"bufio"
	"strings"
	"encoding/csv"
	"io"
	"github.com/jessevdk/go-flags"
	"github.com/gdamore/tcell/v2"
)

func read_stdin() ([]string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func read_files(files []string) ([]string, error) {
	var lines []string

	for _,file := range files {
		fd, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
			return lines, err
		}
		
		scanner := bufio.NewScanner(fd)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		
		fd.Close()
	}
	
	return lines, nil
}

func contains_any(line string, keys []string) (bool) {
	for _,key := range keys {
		if strings.Contains(line, key) {
			return true
		}
	}

	return false
}

func filter_content(lines []string, match_or []string) ([]string) {
	var ret []string
	for i,line := range lines {
		if i>0 && !contains_any(line, match_or) {
			continue
		}

		ret = append(ret, line)
	}
	return ret
}

func tab_parse(lines []string) ([][]string) {
	var res [][]string
	lreader := strings.NewReader(strings.Join(lines, "\n"))
	r := csv.NewReader(lreader)
	for {
		record, err := r.Read()
		if err==io.EOF {
			break
		}

		res = append(res, record)
	}

	return res
}

func calc_limits(tab [][]string) ([]int) {
	ncols := len(tab[0])
	lim := make([]int, ncols)
	for _,row := range(tab) {
		for j,cell := range(row) {
			l := len(cell)
			if(l>lim[j]) {
				lim[j] = l
			}
		}
	}

	return lim
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

func tab_print_plain(tab [][]string) {
	lim := calc_limits(tab)
	
	header_print(tab[0], lim)
	for i,row := range(tab) {
		if i==0 {
			continue
		}
		row_print(row, lim)
	}
}

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

type tviewport struct {
	yoffset int
	ylim_hi int
	screeny int
}

func tcell_line(s tcell.Screen, x, y int, line string) {
	for i, c := range(line) {
		s.SetContent(x+i, y, c, nil, tcell.StyleDefault.Underline(true))
	}
}

func tcell_render(s tcell.Screen, tab [][]string, yoffset int, lim []int) (int,int) {
	screenx, screeny := s.Size()
	s.Clear()

	ylim_lo := 0
	ylim_hi := len(tab) - screeny
	
	if yoffset<0 {
		yoffset = ylim_lo
	}

	if yoffset>ylim_hi {
		yoffset = ylim_hi
	}
	
	for y:=0; y<screeny-1; y++ {
		row := tab[yoffset+y]
		x := 0
		for i,cell := range(row) {
			for k,r := range(cell) {
				xpos := x+k
				if xpos>screenx-1 {
					break
				}
				
				s.SetContent(xpos, y, r, nil, tcell.StyleDefault)
			}
			x += lim[i]
		}
	}

	tcell_line(s, 0, screeny-1, fmt.Sprintf("%d/%d", yoffset, ylim_hi))
	s.Show()

	return yoffset, screeny
}

func tcell_main(tab [][]string) {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	_, screeny := s.Size()
	yoffset := 0
	_,yscreen := s.Size()
	lim := calc_limits(tab)
	
	for {
		yoffset, yscreen = tcell_render(s, tab, yoffset, lim)
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
				yoffset = len(tab) - yscreen
			}
		}
	}
}

func Run() {
	var opts struct {
		Stdin bool `short:"S" long:"stdin" description:"read from stdin"`
		Match []string `short:"m" long:"match" description:"match string"`
		Plain bool `short:"p" long:"plain" description:"dump output, no ui"`
	}

	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	files:=args[1:]
	if len(files)==0 && !opts.Stdin {
		log.Fatal("no input; please supply filenames or enable stdin")
		os.Exit(2)
	}

	var content []string
	if opts.Stdin {
		content, err = read_stdin()
	} else {
		content, err = read_files(files)
	}
	if err != nil || len(content)==0 {
		log.Fatal("no content read")
		os.Exit(3)
	}
	
	if len(opts.Match)>0 {
		content = filter_content(content, opts.Match)
	}

	tab := tab_parse(content)

	if opts.Plain {
		tab_print_plain(tab)
	} else {
		tcell_main(tab)
	}
	os.Exit(0)
}
