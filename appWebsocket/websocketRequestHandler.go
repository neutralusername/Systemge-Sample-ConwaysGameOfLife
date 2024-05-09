package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Websocket"
	"SystemgeSampleApp/typeDefinitions"
)

func (app *App) WebsocketMessageHandler(connection *Websocket.Connection, message *Message.Message) error {
	switch message.TypeName {
	case typeDefinitions.HEARTBEAT_WSREQUEST.Name:
		return nil
	case typeDefinitions.SET_GRID.Name:
		if !typeDefinitions.SET_GRID.Validate(message.Payload) {
			return Error.New("Invalid message payload")
		}
		err := app.messageBroker.Message(typeDefinitions.REQUEST_GRID_CHANGE.New([]string{message.Payload[0][0]}, []string{message.Payload[1][0]}))
		if err != nil {
			return Error.New(err.Error())
		}
		return nil
	default:
		return Error.New("Unknown message type: " + message.TypeName)
	}
}
