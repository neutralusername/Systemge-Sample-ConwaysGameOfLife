package typeDefinitions

import "Systemge/TypeDefinition"

var SET_GRID_WSREQUEST = TypeDefinition.CreateDefinition("set_grid_wsrequest", []int{1, 1}, []string{"row", "col"}, nil)

var GET_GRID_WSPROPAGATE = TypeDefinition.CreateDefinition("get_grid_wspropagate", []int{1}, []string{"grid"}, nil)

var HEARTBEAT = TypeDefinition.CreateDefinition("hb", []int{}, []string{}, nil)
