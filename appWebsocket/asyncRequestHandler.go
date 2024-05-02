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
			return Error.New("Invalid message payload")
		}
		reponse, err := app.gridEndpoint.Request(typeDefinitions.SET_GRID_REQUEST.New([]string{message.Payload[0][0]}, []string{message.Payload[1][0]}))
		if err != nil {
			return Error.New(err.Error())
		}
		app.websocketServer.Broadcast(typeDefinitions.SET_GRID_WSPROPAGATE.New([]string{message.Payload[0][0]}, []string{message.Payload[1][0]}, []string{reponse.Payload[0][0]}))
		return nil
	default:
		return Error.New("Unknown message type: " + message.TypeName)
	}
}
