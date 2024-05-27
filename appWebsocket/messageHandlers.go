package appWebsocket

import (
	"Systemge/Message"
)

func (app *App) WebsocketPropagate(message *Message.Message) error {
	app.websocketServer.Broadcast([]byte(message.Serialize()))
	return nil
}
