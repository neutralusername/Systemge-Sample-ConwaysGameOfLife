package appGameOfLife

import "encoding/json"

type Grid struct {
	Grid [][]int `json:"grid"`
	Rows int     `json:"rows"`
	Cols int     `json:"cols"`
}

func newGrid(grid [][]int, rows, cols int) *Grid {
	return &Grid{
		Grid: grid,
		Rows: rows,
		Cols: cols,
	}
}

func (grid *Grid) marshal() string {
	json, _ := json.Marshal(grid)
	return string(json)
}
