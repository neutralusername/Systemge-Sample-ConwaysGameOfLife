package appWebsocketHTTP

import (
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetSystemgeComponentConfig() *Config.Systemge {
	return &Config.Systemge{
		HandleMessagesSequentially: false,

		SyncRequestTimeoutMs:            10000,
		TcpTimeoutMs:                    5000,
		MaxConnectionAttempts:           0,
		ConnectionAttemptDelayMs:        1000,
		StopAfterOutgoingConnectionLoss: true,
		ServerConfig: &Config.TcpServer{
			Port: 60002,
		},
		EndpointConfigs: []*Config.TcpEndpoint{
			{
				Address: "localhost:60001",
			},
		},
	}
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
