package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Message"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *WebsocketApp) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topic.PROPGATE_GRID:         app.WebsocketPropagate,
		topic.PROPAGATE_GRID_CHANGE: app.WebsocketPropagate,
	}
}
func (app *WebsocketApp) WebsocketPropagate(message *Message.Message) error {
	app.client.GetWebsocketServer().Broadcast(message)
	return nil
}
