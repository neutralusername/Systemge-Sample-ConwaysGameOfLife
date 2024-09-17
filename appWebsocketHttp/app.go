package appWebsocketHttp

import (
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/DashboardClientCustomService"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Metrics"
	"github.com/neutralusername/Systemge/Status"
	"github.com/neutralusername/Systemge/SystemgeConnection"
	"github.com/neutralusername/Systemge/SystemgeMessageHandler"
	"github.com/neutralusername/Systemge/SystemgeServer"
	"github.com/neutralusername/Systemge/WebsocketServer"
)

type AppWebsocketHTTP struct {
	status      int
	statusMutex sync.Mutex

	messageHandlerStopChannel chan<- bool

	systemgeServer  *SystemgeServer.SystemgeServer
	websocketServer *WebsocketServer.WebsocketServer
	httpServer      *HTTPServer.HTTPServer
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{}

	messageHandler := SystemgeMessageHandler.NewTopicExclusiveMessageHandler(
		SystemgeMessageHandler.AsyncMessageHandlers{
			topics.PROPGATE_GRID:         app.websocketPropagate,
			topics.PROPAGATE_GRID_CHANGE: app.websocketPropagate,
		},
		SystemgeMessageHandler.SyncMessageHandlers{},
		nil, nil, 100,
	)
	app.systemgeServer = SystemgeServer.New("systemgeServer",
		&Config.SystemgeServer{
			TcpSystemgeListenerConfig: &Config.TcpSystemgeListener{
				TcpServerConfig: &Config.TcpServer{
					Port: 60001,
				},
			},
			TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
		},
		nil, nil,
		func(connection SystemgeConnection.SystemgeConnection) error {
			stopChannel, _ := SystemgeMessageHandler.StartMessageHandlingLoop_Sequentially(connection, messageHandler)
			app.messageHandlerStopChannel = stopChannel
			return nil
		},
		func(connection SystemgeConnection.SystemgeConnection) {
			close(app.messageHandlerStopChannel)
		},
	)
	app.websocketServer = WebsocketServer.New("websocketServer",
		&Config.WebsocketServer{
			ClientWatchdogTimeoutMs: 1000 * 60,
			Pattern:                 "/ws",
			TcpServerConfig: &Config.TcpServer{
				Port: 8443,
			},
		},
		nil, nil,
		WebsocketServer.MessageHandlers{
			topics.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
			topics.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
			topics.SET_GRID:        app.propagateWebsocketAsyncMessage,
		},
		app.OnConnectHandler, nil,
	)
	app.httpServer = HTTPServer.New("httpServer",
		&Config.HTTPServer{
			TcpServerConfig: &Config.TcpServer{
				Port: 8080,
			},
		},
		nil, nil,
		HTTPServer.Handlers{
			"/": HTTPServer.SendDirectory("../frontend"),
		},
	)

	err := DashboardClientCustomService.New("appWebsocketHttp_dashboardClient",
		&Config.DashboardClient{
			TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
			TcpClientConfig: &Config.TcpClient{
				Address: "localhost:60000",
			},
		}, app,
		nil,
	).Start()
	if err != nil {
		panic(err)
	}
	return app
}

func (app *AppWebsocketHTTP) GetMetrics() map[string]*Metrics.Metrics {
	return map[string]*Metrics.Metrics{}
}

func (app *AppWebsocketHTTP) GetStatus() int {
	return app.status
}

func (app *AppWebsocketHTTP) Start() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STOPPED {
		return Error.New("App already started", nil)
	}
	if err := app.systemgeServer.Start(); err != nil {
		return Error.New("Failed to start systemgeServer", err)
	}
	if err := app.websocketServer.Start(); err != nil {
		app.systemgeServer.Stop()
		return Error.New("Failed to start websocketServer", err)
	}
	if err := app.httpServer.Start(); err != nil {
		app.systemgeServer.Stop()
		app.websocketServer.Stop()
		return Error.New("Failed to start httpServer", err)
	}
	app.status = Status.STARTED
	return nil
}

func (app *AppWebsocketHTTP) Stop() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STARTED {
		return Error.New("App not started", nil)
	}
	app.httpServer.Stop()
	app.websocketServer.Stop()
	app.systemgeServer.Stop()
	app.status = Status.STOPPED
	return nil
}

func (app *AppWebsocketHTTP) websocketPropagate(connection SystemgeConnection.SystemgeConnection, message *Message.Message) {
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
