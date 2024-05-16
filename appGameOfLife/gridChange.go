package appGameOfLife

import "encoding/json"

type GridChange struct {
	row    int
	column int
	state  bool
}

func NewGridChange(row int, column int, state bool) GridChange {
	return GridChange{
		row:    row,
		column: column,
		state:  state,
	}
}

func (gridChange *GridChange) Marshal() string {
	json, _ := json.Marshal(gridChange)
	return string(json)
}

func UnmarshalGridChange(jsonString string) GridChange {
	var gridChange GridChange
	json.Unmarshal([]byte(jsonString), &gridChange)
	return gridChange
}
