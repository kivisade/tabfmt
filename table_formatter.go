package tabfmt

import (
	"unicode/utf8"
	"fmt"
	"strconv"
)

type Table struct {
	rows          []Row
	maxColLengths map[int]int
}

type Row struct {
	cols []string
}

func (tab *Table) AddRow(columns ...string) {
	var row Row
	row.cols = columns
	tab.rows = append(tab.rows, row)
}

func (tab *Table) SetMaxColLengths(maxColLengths map[int]int) {
	tab.maxColLengths = maxColLengths
}

func (tab *Table) getNumCols() int {
	var numCols int = 0
	for _, r := range tab.rows {
		if l := len(r.cols); l > numCols {
			numCols = l
		}
	}
	return numCols
}

func (tab *Table) getColLengths() []int {
	var (
		numCols    int   = tab.getNumCols()
		colLengths []int = make([]int, numCols, numCols)
	)
	for _, r := range tab.rows {
		for j, c := range r.cols {
			if l := utf8.RuneCountInString(c); l > colLengths[j] {
				colLengths[j] = l
			}
		}
	}
	return colLengths
}

func (tab *Table) getCappedColLengths() []int {
	var (
		colLengths []int = tab.getColLengths()
	)
	for i, cl := range colLengths {
		if mcl, found := tab.maxColLengths[i]; found && cl > mcl {
			colLengths[i] = mcl
		}
	}
	return colLengths
}

func (tab *Table) parse() (maxLinesPerRow []int, formattedCell [][][]string) {
	var (
		numRows, numCols int   = len(tab.rows), tab.getNumCols()
		cappedColLengths []int = tab.getCappedColLengths()
		numLines         int
	)

	maxLinesPerRow = make([]int, numRows)
	formattedCell = make([][][]string, numRows)

	for i, r := range tab.rows {
		maxLinesPerRow[i], formattedCell[i] = 0, make([][]string, numCols)
		for j, c := range r.cols {
			formattedCell[i][j] = StringBreakSimple(c, cappedColLengths[j])
			if numLines = len(formattedCell[i][j]); numLines > maxLinesPerRow[i] {
				maxLinesPerRow[i] = numLines
			}
		}
	}

	return
}

func (tab *Table) Print(separator string) {
	var (
		numCols          int   = tab.getNumCols()
		cappedColLengths []int = tab.getCappedColLengths()
		maxLinesPerRow   []int
		formattedCell    [][][]string
		numLines         int
		line             string
	)

	maxLinesPerRow, formattedCell = tab.parse()

	for i, r := range tab.rows {
		for l := 0; l < maxLinesPerRow[i]; l++ {
			for j := range r.cols {
				if numLines = len(formattedCell[i][j]); l < numLines {
					line = formattedCell[i][j][l]
				} else {
					line = ""
				}
				fmt.Printf("%-"+strconv.Itoa(cappedColLengths[j])+"s", line)
				if j < numCols-1 {
					fmt.Print(separator)
				}
			}
			fmt.Println()
		}
	}
}
