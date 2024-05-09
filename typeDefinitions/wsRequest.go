package typeDefinitions

import "Systemge/TypeDefinition"

var HEARTBEAT_WSREQUEST = TypeDefinition.New("heartbeat", []int{}, []string{})
var SET_GRID = TypeDefinition.New("setGrid", []int{1, 1}, []string{"row", "col"})
