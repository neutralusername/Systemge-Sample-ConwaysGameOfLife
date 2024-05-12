package typeDefinitions

import (
	"Systemge/Message"
	"Systemge/TypeDefinition"
	"Systemge/Utilities"
)

var WSPROPAGATE = TypeDefinition.New("wsPropagate", []int{0, 1}, []string{"connectionIds", "message"})

func NewWebsocketMessage(connectionIds []string, message *Message.Message) *Message.Message {
	return WSPROPAGATE.New(connectionIds, []string{Utilities.StringToHexString(string(message.Serialize()))})
}

var GET_GRID = TypeDefinition.New("getGrid", []int{1}, []string{"grid"})
var GET_GRID_CHANGE = TypeDefinition.New("getGridChange", []int{1, 1, 1}, []string{"row", "col", "state"})
