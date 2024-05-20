package appGameOfLife

import "encoding/json"

type Grid struct {
	Grid [GRIDROWS][GRIDCOLS]int `json:"grid"`
	Rows int                     `json:"rows"`
	Cols int                     `json:"cols"`
}

func newGrid(grid [GRIDROWS][GRIDCOLS]int) *Grid {
	return &Grid{
		Grid: grid,
		Rows: GRIDROWS,
		Cols: GRIDCOLS,
	}
}

func (grid *Grid) marshal() string {
	json, _ := json.Marshal(grid)
	return string(json)
}
