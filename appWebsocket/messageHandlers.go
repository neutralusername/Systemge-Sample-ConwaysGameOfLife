package appWebsocket

import (
	"Systemge/Message"
)

func (app *App) GetMessageHandlersSync() map[string]func(*Message.Message) (string, error) {
	return map[string]func(*Message.Message) (string, error){}
}

func (app *App) GetMessageHandlersAsync() map[string]func(*Message.Message) error {
	return map[string]func(*Message.Message) error{
		"getGrid": func(message *Message.Message) error {
			app.websocketServer.Broadcast([]byte(message.Serialize()))
			return nil
		},
		"getGridChange": func(message *Message.Message) error {
			app.websocketServer.Broadcast([]byte(message.Serialize()))
			return nil
		},
	}
}
