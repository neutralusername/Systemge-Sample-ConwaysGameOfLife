package appWebsocket

import (
	"Systemge/ApplicationServer"
	"Systemge/RequestServer"
	"Systemge/WebsocketServer"
)

type App struct {
	appServer       *ApplicationServer.Server
	websocketServer *WebsocketServer.Server

	gridEndpoint RequestServer.Endpoint
}

func Create(appServer *ApplicationServer.Server, websocketServer *WebsocketServer.Server, gridEndpoint RequestServer.Endpoint) WebsocketServer.WebsocketApplication {
	app := &App{
		appServer:       appServer,
		websocketServer: websocketServer,

		gridEndpoint: gridEndpoint,
	}
	return app
}
