package appWebsocketHttp

import (
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"

	"github.com/neutralusername/systemge/configs"
	"github.com/neutralusername/systemge/connectionTcp"
	"github.com/neutralusername/systemge/httpServer"
	"github.com/neutralusername/systemge/listenerWebsocket"
	"github.com/neutralusername/systemge/serviceAccepter"
	"github.com/neutralusername/systemge/serviceTypedReader"
	"github.com/neutralusername/systemge/systemge"
	"github.com/neutralusername/systemge/tools"
)

type AppWebsocketHTTP struct {
	requestResponseManager *tools.RequestResponseManager[*tools.Message]
	internalConnection     systemge.Connection[[]byte]
	websocketConnections   map[systemge.Connection[[]byte]]struct{}
	mutex                  sync.RWMutex
}

func New() *AppWebsocketHTTP {
	internalConnection, err := connectionTcp.EstablishConnection(
		&configs.TcpBufferedReader{},
		&configs.TcpClient{
			Address: "localhost:60001",
		},
	)
	if err != nil {
		panic(err)
	}
	app := &AppWebsocketHTTP{
		requestResponseManager: tools.NewRequestResponseManager[*tools.Message](&configs.RequestResponseManager{}),
		internalConnection:     internalConnection,
		websocketConnections:   make(map[systemge.Connection[[]byte]]struct{}),
	}

	reader, err := serviceTypedReader.NewAsync(
		internalConnection,
		&configs.ReaderAsync{},
		&configs.Routine{
			MaxConcurrentHandlers: 10,
		},
		func(message *tools.Message, connection systemge.Connection[[]byte]) {
			if message.GetSyncToken() != "" {
				if !message.IsResponse() {
					panic("message is not a response")
				}
				if err := app.requestResponseManager.AddResponse(message.GetSyncToken(), message); err != nil {
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
							message.Serialize(),
							0,
						); err != nil {
							panic(err)
						}
					}()
				}
			}
		},
		func(data []byte) (*tools.Message, error) {
			return tools.DeserializeMessage(data)
		},
	)
	if err != nil {
		panic(err)
	}
	if err := reader.GetRoutine().Start(); err != nil {
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
	if err := httpServer.Start(); err != nil {
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
	if err := listenerWebsocket.Start(); err != nil {
		panic(err)
	}

	websocketAccepter, err := serviceAccepter.New(
		listenerWebsocket,
		&configs.Accepter{},
		&configs.Routine{
			MaxConcurrentHandlers: 1,
		},
		func(connection systemge.Connection[[]byte]) error {
			reader, err := serviceTypedReader.NewAsync(
				connection,
				&configs.ReaderAsync{},
				&configs.Routine{
					MaxConcurrentHandlers: 10,
				},
				func(message *tools.Message, connection systemge.Connection[[]byte]) {
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
							message.Serialize(),
							0,
						); err != nil {
							panic(err)
						}
					}()
				},
				func(data []byte) (*tools.Message, error) {
					return tools.DeserializeMessage(data)
				},
			)
			if err != nil {
				panic(err)
			}
			if err := reader.GetRoutine().Start(); err != nil {
				panic(err)
			}

			app.mutex.Lock()
			app.websocketConnections[connection] = struct{}{}
			app.mutex.Unlock()

			go func() { // abstract on close handler
				<-connection.GetCloseChannel()

				app.mutex.Lock()
				delete(app.websocketConnections, connection)
				app.mutex.Unlock()
			}()

			request, err := app.requestResponseManager.NewRequest(tools.GenerateRandomString(32, tools.ALPHA_NUMERIC), 1, 0, nil)
			if err != nil {
				panic(err)
			}

			if err := app.internalConnection.Write(
				tools.NewMessage(
					topics.GET_GRID,
					"",
					request.GetToken(),
					false,
				).Serialize(),
				0,
			); err != nil {
				panic(err)
			}

			response, err := request.GetNextResponse()
			if err != nil {
				panic(err)
			}

			if err = connection.Write(
				response.Serialize(),
				0,
			); err != nil {
				panic(err)
			}
			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	if err = websocketAccepter.GetRoutine().Start(); err != nil {
		panic(err)
	}

	return app
}
