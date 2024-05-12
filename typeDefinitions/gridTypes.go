package typeDefinitions

import (
	"Systemge/TypeDefinition"
)

var REQUEST_GRID = TypeDefinition.New("requestGridUnicast", []int{0}, []string{"connectionIds"})
var REQUEST_GRID_CHANGE = TypeDefinition.New("requestGridChange", []int{1, 1}, []string{"row", "col"})

var GET_GRID = TypeDefinition.New("getGrid", []int{1}, []string{"grid"})
var GET_GRID_CHANGE = TypeDefinition.New("getGridChange", []int{1, 1, 1}, []string{"row", "col", "state"})
