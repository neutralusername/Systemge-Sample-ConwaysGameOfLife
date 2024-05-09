package appWebsocket

import (
	"Systemge/Message"
	"SystemgeSampleApp/typeDefinitions"
	"errors"
)

func (app *App) MessageHandler(message *Message.Message) error {
	switch message.TypeName {
	case typeDefinitions.BROADCAST_GRID.Name:
		app.websocketServer.Broadcast(typeDefinitions.GET_GRID.New([]string{message.Payload[0][0]}))
	case typeDefinitions.BROADCAST_GRID_CHANGE.Name:
		app.websocketServer.Broadcast(typeDefinitions.GET_GRID_CHANGE.New([]string{message.Payload[0][0]}, []string{message.Payload[1][0]}, []string{message.Payload[2][0]}))
	case typeDefinitions.UNICAST_GRID.Name:
		app.websocketServer.Unicast(message.Payload[0][0], typeDefinitions.GET_GRID.New([]string{message.Payload[1][0]}))
	default:
		return errors.New("Unknown message type: " + message.TypeName)
	}
	return nil
}
