package appWebsocketHttp

import (
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"

	"github.com/neutralusername/systemge/configs"
	"github.com/neutralusername/systemge/httpServer"
	"github.com/neutralusername/systemge/listenerWebsocket"
	"github.com/neutralusername/systemge/server"
	"github.com/neutralusername/systemge/serviceReader"
	"github.com/neutralusername/systemge/systemge"
	"github.com/neutralusername/systemge/tools"
	"github.com/neutralusername/systemge/typedListener"
)

type AppWebsocketHTTP struct {
	requestResponseManager *tools.RequestResponseManager[*tools.Message]
	internalConnection     systemge.Connection[*tools.Message]
	websocketConnections   map[systemge.Connection[*tools.Message]]struct{}
	mutex                  sync.RWMutex
}

func New() *AppWebsocketHTTP {
	internalConnection, err := appGameOfLife.Connector.Connect(0)
	if err != nil {
		panic(err)
	}

	app := &AppWebsocketHTTP{
		requestResponseManager: tools.NewRequestResponseManager[*tools.Message](&configs.RequestResponseManager{}),
		internalConnection:     internalConnection,
		websocketConnections:   make(map[systemge.Connection[*tools.Message]]struct{}),
	}

	internalConnectionReader, err := serviceReader.NewAsync(
		app.internalConnection,
		&configs.ReaderAsync{},
		&configs.Routine{
			MaxConcurrentHandlers: 10,
		},
		app.internalReadHandler,
	)
	if err != nil {
		panic(err)
	}
	if err := internalConnectionReader.GetRoutine().Start(); err != nil {
		panic(err)
	}

	httpServer, err := httpServer.New(
		"httpServer",
		&configs.HTTPServer{
			TcpListenerConfig: &configs.TcpListener{
				Port:   8080,
				Domain: "localhost",
			},
		},
		nil,
		httpServer.HandlerFuncs{
			"/": httpServer.SendDirectory("../frontend"),
		},
	)
	if err != nil {
		panic(err)
	}
	if err := httpServer.Start(); err != nil {
		panic(err)
	}

	listenerWebsocket, err := listenerWebsocket.New(
		"listenerWebsocket",
		nil,
		&configs.WebsocketListener{
			TcpListenerConfig: &configs.TcpListener{
				Port: 8443,
			},
			Pattern: "/ws",
		},
		0,
		0,
	)
	if err != nil {
		panic(err)
	}
	if err := listenerWebsocket.Start(); err != nil {
		panic(err)
	}

	typedListenerWebsocket, err := typedListener.New(
		listenerWebsocket,
		tools.DeserializeMessage,
		tools.SerializeMessage,
	)
	if err != nil {
		panic(err)
	}

	server, err := server.New(
		typedListenerWebsocket,
		&configs.Accepter{},
		&configs.Routine{
			MaxConcurrentHandlers: 1,
		},
		app.websocketAcceptHandler,
		&configs.ReaderAsync{},
		&configs.Routine{
			MaxConcurrentHandlers: 1,
		},
		app.websocketReadHandler,
	)
	server.GetAccepter().GetRoutine().Start()

	return app
}

func (app *AppWebsocketHTTP) internalReadHandler(message *tools.Message, connection systemge.Connection[*tools.Message]) {
	if message.GetSyncToken() != "" {
		if !message.IsResponse() {
			panic("message is not a response")
		}
		if err := app.requestResponseManager.AddResponse(
			message.GetSyncToken(),
			message,
		); err != nil {
			panic(err)
		}
	} else {
		switch message.GetTopic() {
		case topics.PROPAGATE_GRID:
		case topics.PROPAGATE_GRID_CHANGE:
		default:
			panic("unknown topic")
		}

		for websocketConnection := range app.websocketConnections {
			go func() {
				if err := websocketConnection.Write(
					message,
					0,
				); err != nil {
					panic(err)
				}
			}()
		}
	}
}

func (app *AppWebsocketHTTP) websocketAcceptHandler(connection systemge.Connection[*tools.Message]) error {
	app.mutex.Lock()
	app.websocketConnections[connection] = struct{}{}
	app.mutex.Unlock()

	go func() { // abstract on close handler
		<-connection.GetCloseChannel()

		app.mutex.Lock()
		delete(app.websocketConnections, connection)
		app.mutex.Unlock()
	}()

	request, err := app.requestResponseManager.NewRequest(
		tools.GenerateRandomString(32, tools.ALPHA_NUMERIC),
		1,
		0,
		nil,
	)
	if err != nil {
		panic(err)
	}

	if err := app.internalConnection.Write(
		tools.NewMessage(
			topics.GET_GRID,
			"",
			request.GetToken(),
			false,
		),
		0,
	); err != nil {
		panic(err)
	}

	response, err := request.GetNextResponse()
	if err != nil {
		panic(err)
	}

	if err = connection.Write(
		response,
		0,
	); err != nil {
		panic(err)
	}
	return nil
}

func (app *AppWebsocketHTTP) websocketReadHandler(message *tools.Message, connection systemge.Connection[*tools.Message]) {
	switch message.GetTopic() {
	case topics.GRID_CHANGE:
	case topics.NEXT_GENERATION:
	case topics.SET_GRID:
		/* 	case "heartbeat":
		return */
	default:
		panic("unknown topic")
	}

	go func() {
		if err := app.internalConnection.Write(
			message,
			0,
		); err != nil {
			panic(err)
		}
	}()
}
