package WebsocketApp

import (
	"Systemge/Error"
	"Systemge/Message"
	"strings"
)

func (app *App) MessageHandler(message *Message.Message) error {
	switch message.Type {
	case "websocketUnicast":
		segments := strings.Split(message.Body, "|")
		app.websocketServer.Unicast(segments[0], []byte(Message.New(segments[1], app.name, segments[2]).Serialize()))
		return nil
	case "getGrid":
		app.websocketServer.Broadcast([]byte(message.Serialize()))
		return nil
	case "getGridChange":
		app.websocketServer.Broadcast([]byte(message.Serialize()))
		return nil
	default:
		return Error.New("Unknown message type")
	}
}
