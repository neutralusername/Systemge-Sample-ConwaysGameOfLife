package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleApp/typeDefinitions"
)

func (app *WebsocketApp) MessageHandler(message *Message.Message) error {
	switch message.TypeName {
	case typeDefinitions.WSPROPAGATE.Name:
		msg := Message.Deserialize(Utilities.HexStringToString(message.Payload[1][0]))
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
