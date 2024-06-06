package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Message"
	"SystemgeSampleApp/topic"
)

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topic.GET_GRID:        app.WebsocketPropagate,
		topic.GET_GRID_CHANGE: app.WebsocketPropagate,
	}
}
func (app *App) WebsocketPropagate(message *Message.Message) error {
	app.messageBrokerClient.GetWebsocketServer().Broadcast([]byte(message.Serialize()))
	return nil
}
