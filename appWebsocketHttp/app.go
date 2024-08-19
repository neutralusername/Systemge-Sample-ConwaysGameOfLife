package appWebsocketHttp

import (
	"SystemgeSampleConwaysGameOfLife/topics"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeMessageHandler"
	"github.com/neutralusername/Systemge/SystemgeServer"
	"github.com/neutralusername/Systemge/WebsocketServer"
)

type AppWebsocketHTTP struct {
	systemgeServer  *SystemgeServer.SystemgeServer
	websocketServer *WebsocketServer.WebsocketServer
	httpServer      *HTTPServer.HTTPServer
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{}
	app.systemgeServer = SystemgeServer.New(&Config.SystemgeServer{
		Name: "systemgeServer",
		ListenerConfig: &Config.SystemgeListener{
			TcpListenerConfig: &Config.TcpListener{
				Port: 60001,
			},
		},
		ConnectionConfig: &Config.SystemgeConnection{},
	}, &Config.SystemgeReceiver{},
		SystemgeMessageHandler.New(SystemgeMessageHandler.AsyncMessageHandlers{
			topics.PROPGATE_GRID:         app.WebsocketPropagate,
			topics.PROPAGATE_GRID_CHANGE: app.WebsocketPropagate,
		}, SystemgeMessageHandler.SyncMessageHandlers{}))
	app.websocketServer = WebsocketServer.New(&Config.WebsocketServer{
		ClientWatchdogTimeoutMs: 1000 * 60,
		Pattern:                 "/ws",
		TcpListenerConfig: &Config.TcpListener{
			Port: 8443,
		},
	}, map[string]WebsocketServer.MessageHandler{
		topics.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
		topics.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
		topics.SET_GRID:        app.propagateWebsocketAsyncMessage,
	}, app.OnConnectHandler, nil)
	app.httpServer = HTTPServer.New(&Config.HTTPServer{
		TcpListenerConfig: &Config.TcpListener{
			Port: 8080,
		},
	}, HTTPServer.Handlers{
		"/": HTTPServer.SendDirectory("../frontend"),
	})
	if err := app.systemgeServer.Start(); err != nil {
		panic(err)
	}
	app.websocketServer.Start()
	app.httpServer.Start()
	return app
}

func (app *AppWebsocketHTTP) WebsocketPropagate(message *Message.Message) {
	app.websocketServer.Broadcast(message)
}

func (app *AppWebsocketHTTP) OnConnectHandler(websocketClient *WebsocketServer.WebsocketClient) error {
	responseChannel, err := app.systemgeServer.SyncRequest(topics.GET_GRID, websocketClient.GetId())
	if err != nil {
		return Error.New("Failed to get grid", err)
	}
	response := <-responseChannel
	if response == nil {
		return Error.New("Failed to receive response", err)
	}
	getGridMessage := Message.NewAsync(topics.GET_GRID, response.GetPayload())
	websocketClient.Send([]byte(getGridMessage.Serialize()))
	return nil
}

func (app *AppWebsocketHTTP) propagateWebsocketAsyncMessage(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
	return app.systemgeServer.AsyncMessage(message.GetTopic(), message.GetPayload())
}
