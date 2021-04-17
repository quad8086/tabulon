package tabulon

import (
	"fmt"
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
