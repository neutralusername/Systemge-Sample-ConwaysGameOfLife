package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
)

func (app *AppWebsocketHTTP) GetSystemgeComponentConfig() *Config.Systemge {
	return app.systemgeConfig
}

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return app.asyncMessageHandlers
}
func (app *AppWebsocketHTTP) WebsocketPropagate(node *Node.Node, message *Message.Message) error {
	node.WebsocketBroadcast(message)
	return nil
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return app.syncMessageHandlers
}
