package appGameOfLife

import "encoding/json"

type GridChange struct {
	Row    int `json:"row"`
	Column int `json:"column"`
	State  int `json:"state"`
}

func NewGridChange(row int, column int, state int) GridChange {
	return GridChange{
		Row:    row,
		Column: column,
		State:  state,
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
