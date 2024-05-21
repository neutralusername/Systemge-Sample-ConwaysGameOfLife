package appWebsocket

import (
	"Systemge/Message"
)

func (app *App) GetGrid(message *Message.Message) error {
	app.websocketServer.Broadcast([]byte(message.Serialize()))
	return nil
}

func (app *App) GetGridChange(message *Message.Message) error {
	app.websocketServer.Broadcast([]byte(message.Serialize()))
	return nil
}
