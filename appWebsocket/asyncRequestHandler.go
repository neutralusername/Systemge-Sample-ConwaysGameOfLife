package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/WebsocketServer"
	"SystemgeSampleApp/typeDefinitions"
)

func (app *App) AsyncRequestHandler(connection *WebsocketServer.Connection, message *Message.Message) error {
	switch message.TypeName {
	case typeDefinitions.HEARTBEAT.Name:
		return nil
	case typeDefinitions.SET_GRID_WSREQUEST.Name:
		if !typeDefinitions.SET_GRID_WSREQUEST.Validate(message.Payload) {
			return Error.Create("Invalid message payload")
		}
		reponse, err := app.gridEndpoint.Request(typeDefinitions.SET_GRID_REQUEST.Create([]string{message.Payload[0][0]}, []string{message.Payload[1][0]}))
		if err != nil {
			return Error.Create(err.Error())
		}
		app.websocketServer.Broadcast(typeDefinitions.GET_GRID_WSPROPAGATE.Create([]string{reponse.Payload[0][0]}))
		return nil
	default:
		return Error.Create("Unknown message type: " + message.TypeName)
	}
}
