package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleConwaysGameOfLife/topics"
)

func (app *AppWebsocketHTTP) GetSystemgeComponentConfig() *Config.Systemge {
	return &Config.Systemge{
		HandleMessagesSequentially: false,

		BrokerSubscribeDelayMs:    1000,
		TopicResolutionLifetimeMs: 10000,
		SyncResponseTimeoutMs:     10000,
		TcpTimeoutMs:              5000,

		ResolverEndpoint: &Config.TcpEndpoint{
			Address: "127.0.0.1:60000",
		},
	}
}

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.PROPGATE_GRID:         app.WebsocketPropagate,
		topics.PROPAGATE_GRID_CHANGE: app.WebsocketPropagate,
	}
}
func (app *AppWebsocketHTTP) WebsocketPropagate(node *Node.Node, message *Message.Message) error {
	node.WebsocketBroadcast(message)
	return nil
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}
