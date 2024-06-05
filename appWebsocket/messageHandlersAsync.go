package appWebsocket

import (
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"SystemgeSampleApp/topic"
)

func (app *App) GetAsyncMessageHandlers() map[string]MessageBrokerClient.AsyncMessageHandler {
	return map[string]MessageBrokerClient.AsyncMessageHandler{
		topic.GET_GRID:        app.WebsocketPropagate,
		topic.GET_GRID_CHANGE: app.WebsocketPropagate,
	}
}
func (app *App) WebsocketPropagate(message *Message.Message) error {
	app.messageBrokerClient.WebsocketBroadcast([]byte(message.Serialize()))
	return nil
}
