package typeDefinitions

import (
	"Systemge/Message"
	"Systemge/TypeDefinition"
)

var WSPROPAGATE_MESSAGE_TYPE = TypeDefinition.New("wsPropagate", []int{0, 1, 0}, []string{"connectionIds", "messageType", "paylods"})

func NewWebsocketMessage(connectionIds []string, messageType string, payloads []string) *Message.Message {
	return WSPROPAGATE_MESSAGE_TYPE.New(connectionIds, []string{messageType}, payloads)
}
