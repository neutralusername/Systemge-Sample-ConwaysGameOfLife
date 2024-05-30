package appWebsocket

import (
	"Systemge/Message"
)

func (app *App) WebsocketPropagate(message *Message.Message) error {
	app.messageBrokerClient.WebsocketBroadcast([]byte(message.Serialize()))
	return nil
}
