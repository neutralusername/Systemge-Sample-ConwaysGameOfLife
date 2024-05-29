package dto

import "encoding/json"

type Grid struct {
	Grid [][]int `json:"grid"`
	Rows int     `json:"rows"`
	Cols int     `json:"cols"`
}

func NewGrid(grid [][]int, rows, cols int) *Grid {
	return &Grid{
		Grid: grid,
		Rows: rows,
		Cols: cols,
	}
}

func (grid *Grid) Marshal() string {
	json, _ := json.Marshal(grid)
	return string(json)
}
