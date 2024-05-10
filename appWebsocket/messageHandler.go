package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"SystemgeSampleApp/typeDefinitions"
)

func (app *WebsocketApp) MessageHandler(message *Message.Message) error {
	switch message.TypeName {
	case typeDefinitions.WSPROPAGATE_MESSAGE_TYPE.Name:
		msg := &Message.Message{
			TypeName: message.Payload[1][0],
			Payload:  [][]string{message.Payload[2]},
		}
		if len(message.Payload[0]) == 0 {
			app.websocketServer.Broadcast(msg.Serialize())
		} else {
			app.websocketServer.Multicast(message.Payload[0], msg.Serialize())
		}
	default:
		return Error.New("Unknown message type: " + message.TypeName)
	}
	return nil
}
