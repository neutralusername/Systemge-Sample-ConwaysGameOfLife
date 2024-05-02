package typeDefinitions

import "Systemge/TypeDefinition"

var SET_GRID_WSREQUEST = TypeDefinition.New("set_grid_wsrequest", []int{1, 1}, []string{"row", "col"}, nil)

var GET_GRID_WSPROPAGATE = TypeDefinition.New("get_grid_wspropagate", []int{1}, []string{"grid"}, nil)

var SET_GRID_WSPROPAGATE = TypeDefinition.New("set_grid_wspropagate", []int{1, 1, 1}, []string{"row", "col", "state"}, nil)

var HEARTBEAT = TypeDefinition.New("hb", []int{}, []string{}, nil)
