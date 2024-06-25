package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Message"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Client.AsyncMessageHandler {
	return map[string]Client.AsyncMessageHandler{
		topic.PROPGATE_GRID:         app.WebsocketPropagate,
		topic.PROPAGATE_GRID_CHANGE: app.WebsocketPropagate,
	}
}
func (app *AppWebsocketHTTP) WebsocketPropagate(client *Client.Client, message *Message.Message) error {
	client.WebsocketBroadcast(message)
	return nil
}
