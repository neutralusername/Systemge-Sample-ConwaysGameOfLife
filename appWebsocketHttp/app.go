package appWebsocketHttp

import (
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"

	"github.com/neutralusername/systemge/configs"
	"github.com/neutralusername/systemge/connectionChannel"
	"github.com/neutralusername/systemge/httpServer"
	"github.com/neutralusername/systemge/listenerWebsocket"
	"github.com/neutralusername/systemge/serviceAccepter"
	"github.com/neutralusername/systemge/serviceReader"
	"github.com/neutralusername/systemge/serviceTypedReader"
	"github.com/neutralusername/systemge/systemge"
	"github.com/neutralusername/systemge/tools"
)

type AppWebsocketHTTP struct {
	requestResponseManager *tools.RequestResponseManager[*tools.Message]

	websocketAccepter *serviceAccepter.Accepter[[]byte]

	listenerWebsocket systemge.Listener[[]byte, systemge.Connection[[]byte]]
	httpServer        *httpServer.HTTPServer

	internalConnection       systemge.Connection[*tools.Message]
	internalConnectionReader *serviceReader.ReaderAsync[*tools.Message]

	websocketConnections map[systemge.Connection[[]byte]]struct{}
	mutex                sync.RWMutex
}

func New() *AppWebsocketHTTP {

	connChan := appGameOfLife.ConnectionChannel
	if connChan == nil {
		panic("connection channel is nil")
	}

	internalConnection, err := connectionChannel.EstablishConnection(
		connChan,
		0,
	)
	if err != nil {
		panic(err)
	}
	app := &AppWebsocketHTTP{
		requestResponseManager: tools.NewRequestResponseManager[*tools.Message](&configs.RequestResponseManager{}),
		internalConnection:     internalConnection,
	}

	reader, err := serviceReader.NewAsync(
		internalConnection,
		&configs.ReaderAsync{},
		&configs.Routine{
			MaxConcurrentHandlers: 1,
		},
		func(message *tools.Message, connection systemge.Connection[*tools.Message]) {

			if message.GetSyncToken() != "" {
				if !message.IsResponse() {
					return
				}
				if err := app.requestResponseManager.AddResponse(message.GetSyncToken(), message); err != nil {
					return
				}
			} else {
				switch message.GetTopic() {
				case topics.PROPAGATE_GRID:
				case topics.PROPAGATE_GRID_CHANGE:
				default:
					return
				}

				app.mutex.RLock()
				defer app.mutex.RUnlock()

				for websocketConnection := range app.websocketConnections {
					go websocketConnection.Write(message.Serialize(), 0)
				}
			}
		},
	)
	if err != nil {
		panic(err)
	}
	app.internalConnectionReader = reader
	err = app.internalConnectionReader.GetRoutine().Start()
	if err != nil {
		panic(err)
	}

	httpServer, err := httpServer.New(
		"httpServer",
		&configs.HTTPServer{
			TcpServerConfig: &configs.TcpServer{
				Port: 8080,
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
	app.httpServer = httpServer
	err = app.httpServer.Start()
	if err != nil {
		panic(err)
	}

	listenerWebsocket, err := listenerWebsocket.New(
		"listenerWebsocket",
		nil,
		&configs.WebsocketListener{
			TcpServerConfig: &configs.TcpServer{
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
	app.listenerWebsocket = listenerWebsocket
	err = app.listenerWebsocket.Start()
	if err != nil {
		panic(err)
	}

	websocketAccepter, err := serviceAccepter.New(
		listenerWebsocket,
		&configs.Accepter{},
		&configs.Routine{
			MaxConcurrentHandlers: 1,
		},
		func(connection systemge.Connection[[]byte]) error {

			_, err := serviceTypedReader.NewAsync(
				connection,
				&configs.ReaderAsync{},
				&configs.Routine{},
				func(message *tools.Message, connection systemge.Connection[[]byte]) {

					switch message.GetTopic() {
					case topics.GRID_CHANGE:
					case topics.NEXT_GENERATION:
					case topics.SET_GRID:
					default:
						return
					}

					app.mutex.RLock()
					defer app.mutex.RUnlock()

					go app.internalConnection.Write(message, 0)
				},
				func(data []byte) (*tools.Message, error) {
					return tools.DeserializeMessage(data)
				},
			)
			if err != nil {
				return err
			}

			app.mutex.Lock()
			app.websocketConnections[connection] = struct{}{}
			app.mutex.Unlock()

			go func() { // abstract on close handler
				<-connection.GetCloseChannel()
				app.mutex.Lock()
				defer app.mutex.Unlock()

				delete(app.websocketConnections, connection)
			}()

			request, err := app.requestResponseManager.NewRequest(tools.GenerateRandomString(32, tools.ALPHA_NUMERIC), 1, 0, nil)
			if err != nil {
				panic(err)
			}

			if err = app.internalConnection.Write(tools.NewMessage(topics.GET_GRID, "", request.GetToken(), false), 0); err != nil {
				panic(err)
			}

			response, err := request.GetNextResponse()
			if err != nil {
				panic(err)
			}

			if err = connection.Write(response.Serialize(), 0); err != nil {
				panic(err)
			}

			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	app.websocketAccepter = websocketAccepter
	err = app.websocketAccepter.GetRoutine().Start()
	if err != nil {
		panic(err)
	}

	return app
}
