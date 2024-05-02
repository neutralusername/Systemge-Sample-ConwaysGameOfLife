package appWebsocket

import (
	"Systemge/Message"
	"Systemge/TypeDefinition"
	"SystemgeSampleApp/typeDefinitions"
	"errors"
)

func (app *App) SyncRequestHandler(message *Message.Message) *Message.Message {
	switch message.TypeName {
	case typeDefinitions.PROPAGATE_GRID_REQUEST.Name:
		app.websocketServer.Broadcast(typeDefinitions.GET_GRID_WSPROPAGATE.New([]string{message.Payload[0][0]}))
		return typeDefinitions.PROPAGATE_GRID_REQUEST.Response.New()
	default:
		return TypeDefinition.NewErrorMessage(errors.New("Unknown message type: " + message.TypeName))
	}
}
