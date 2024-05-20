package appGameOfLife

import "encoding/json"

type Grid struct {
	Grid [GRIDROWS][GRIDCOLS]int `json:"grid"`
	Rows int                     `json:"rows"`
	Cols int                     `json:"cols"`
}

func NewGrid(grid [GRIDROWS][GRIDCOLS]int) *Grid {
	return &Grid{
		Grid: grid,
		Rows: GRIDROWS,
		Cols: GRIDCOLS,
	}
}

func (grid *Grid) Marshal() string {
	json, _ := json.Marshal(grid)
	return string(json)
}

func UnmarshalGrid(jsonString string) Grid {
	var grid Grid
	json.Unmarshal([]byte(jsonString), &grid)
	return grid
}
