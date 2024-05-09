package typeDefinitions

import "Systemge/TypeDefinition"

var BROADCAST_GRID = TypeDefinition.New("broadcastGrid", []int{1}, []string{"grid"})
var BROADCAST_GRID_CHANGE = TypeDefinition.New("broadcastGridChange", []int{1, 1, 1}, []string{"row", "col", "state"})
var UNICAST_GRID = TypeDefinition.New("unicastGrid", []int{1, 1}, []string{"connectionId", "grid"})
