package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Node"
	"SystemgeSampleConwaysGameOfLife/topics"
	"net/http"

	"github.com/gorilla/websocket"
)

type AppWebsocketHTTP struct {
	websocketMessageHandlers map[string]Node.WebsocketMessageHandler
	syncMessageHandlers      map[string]Node.SyncMessageHandler
	asyncMessageHandlers     map[string]Node.AsyncMessageHandler
	systemgeConfig           *Config.Systemge
	websocketConfig          *Config.Websocket
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{}
	app.websocketMessageHandlers = map[string]Node.WebsocketMessageHandler{
		topics.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
		topics.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
		topics.SET_GRID:        app.propagateWebsocketAsyncMessage,
	}
	app.syncMessageHandlers = map[string]Node.SyncMessageHandler{}
	app.asyncMessageHandlers = map[string]Node.AsyncMessageHandler{
		topics.PROPGATE_GRID:         app.WebsocketPropagate,
		topics.PROPAGATE_GRID_CHANGE: app.WebsocketPropagate,
	}
	app.systemgeConfig = &Config.Systemge{
		HandleMessagesSequentially: false,

		BrokerSubscribeDelayMs:    1000,
		TopicResolutionLifetimeMs: 10000,
		SyncResponseTimeoutMs:     10000,
		TcpTimeoutMs:              5000,

		ResolverEndpoint: &Config.TcpEndpoint{
			Address: "127.0.0.1:60000",
		},
	}
	app.websocketConfig = &Config.Websocket{
		Pattern: "/ws",
		Server: &Config.TcpServer{
			Port:      8443,
			Blacklist: []string{},
			Whitelist: []string{},
		},
		HandleClientMessagesSequentially: false,
		ClientMessageCooldownMs:          0,
		ClientWatchdogTimeoutMs:          20000,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	return app
}
