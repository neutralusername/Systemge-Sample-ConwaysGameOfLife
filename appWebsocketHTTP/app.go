package appWebsocketHTTP

import (
	"SystemgeSampleConwaysGameOfLife/topics"

	"github.com/neutralusername/Systemge/Node"
)

type AppWebsocketHTTP struct {
	websocketMessageHandlers map[string]Node.WebsocketMessageHandler
	syncMessageHandlers      map[string]Node.SyncMessageHandler
	asyncMessageHandlers     map[string]Node.AsyncMessageHandler
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

	return app
}
