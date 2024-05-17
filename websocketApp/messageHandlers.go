package WebsocketApp

import (
	"Systemge/Message"
	"strings"
)

func (app *App) WebsocketUnicast(message *Message.Message) error {
	segments := strings.Split(message.Body, "|")
	app.websocketServer.Unicast(segments[0], []byte(Message.New(segments[1], app.name, segments[2]).Serialize()))
	return nil
}

func (app *App) GetGrid(message *Message.Message) error {
	app.websocketServer.Broadcast([]byte(message.Serialize()))
	return nil
}

func (app *App) GetGridChange(message *Message.Message) error {
	app.websocketServer.Broadcast([]byte(message.Serialize()))
	return nil
}
