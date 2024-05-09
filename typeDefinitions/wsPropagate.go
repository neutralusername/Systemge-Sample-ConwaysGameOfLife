package typeDefinitions

import "Systemge/TypeDefinition"

var GET_GRID = TypeDefinition.New("getGrid", []int{1}, []string{"grid"})
var GET_GRID_CHANGE = TypeDefinition.New("getGridChange", []int{1, 1, 1}, []string{"row", "col", "state"})
