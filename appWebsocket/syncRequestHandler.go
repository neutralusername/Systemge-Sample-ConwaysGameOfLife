package appWebsocket

import (
	"Systemge/Message"
	"Systemge/TypeDefinition"
	"errors"
)

func (app *App) SyncRequestHandler(message *Message.Message) *Message.Message {
	switch message.TypeName {
	default:
		return TypeDefinition.CreateError(errors.New("Unknown message type: " + message.TypeName))
	}
}
