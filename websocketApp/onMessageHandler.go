package WebsocketApp

import (
	"Systemge/Message"
	"Systemge/Websocket"
)

func (app *App) OnMessageHandler(connection *Websocket.Connection, message *Message.Message) {
	switch message.Type {
	case "heartbeat":
		connection.ResetWatchdog()
	default:
		err := app.messageBrokerClient.Send(message)
		if err != nil {
			connection.SendMessage([]byte(Message.New("error", app.name, err.Error()).Serialize()))
		}
	}
}
