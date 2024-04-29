package typeDefinitions

import (
	"Systemge/TypeDefinition"
)

var SET_GRID_REQUEST = TypeDefinition.CreateDefinition("set_grid_request", []int{1, 1}, []string{"row", "col"}, &SET_GRID_RESPONSE)
var SET_GRID_RESPONSE = TypeDefinition.CreateDefinition("set_grid_response", []int{0}, []string{"grid"}, nil)

var GET_GRID_REQUEST = TypeDefinition.CreateDefinition("get_grid_request", []int{}, []string{}, &GET_GRID_RESPONSE)
var GET_GRID_RESPONSE = TypeDefinition.CreateDefinition("get_grid_response", []int{0}, []string{"grid"}, nil)
