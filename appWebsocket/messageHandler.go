package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/TypeDefinition"
)

func (app *WebsocketApp) MessageHandler(message *Message.Message) error {
	switch message.TypeName {
	case TypeDefinition.WSPROPAGATE_MESSAGE_TYPE.Name:
		app.websocketServer.PropagateMessage(message.Payload[0], message.Payload[1][0], message.Payload[2])
	default:
		return Error.New("Unknown message type: " + message.TypeName)
	}
	return nil
}
