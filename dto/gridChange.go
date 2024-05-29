package dto

import "encoding/json"

type GridChange struct {
	Row    int `json:"row"`
	Column int `json:"column"`
	State  int `json:"state"`
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
