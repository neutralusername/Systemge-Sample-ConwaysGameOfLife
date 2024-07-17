package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *AppWebsocketHTTP) GetSystemgeConfig() Config.Systemge {
	return Config.Systemge{
		HandleMessagesSequentially: false,
	}
}

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topic.PROPGATE_GRID:         app.WebsocketPropagate,
		topic.PROPAGATE_GRID_CHANGE: app.WebsocketPropagate,
	}
}
func (app *AppWebsocketHTTP) WebsocketPropagate(node *Node.Node, message *Message.Message) error {
	node.WebsocketBroadcast(message)
	return nil
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}
