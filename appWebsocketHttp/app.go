package appWebsocketHttp

import (
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"

	"github.com/neutralusername/systemge/configs"
	"github.com/neutralusername/systemge/httpServer"
	"github.com/neutralusername/systemge/listenerChannel"
	"github.com/neutralusername/systemge/listenerWebsocket"
	"github.com/neutralusername/systemge/serviceAccepter"
	"github.com/neutralusername/systemge/serviceReader"
	"github.com/neutralusername/systemge/serviceTypedReader"
	"github.com/neutralusername/systemge/systemge"
	"github.com/neutralusername/systemge/tools"
)

type AppWebsocketHTTP struct {
	requestResponseManager *tools.RequestResponseManager[*tools.Message]

	channelAccepter   *serviceAccepter.Accepter[*tools.Message]
	websocketAccepter *serviceAccepter.Accepter[[]byte]

	websocketConnections map[systemge.Connection[[]byte]]struct{}
	internalConnections  map[systemge.Connection[*tools.Message]]struct{}
	mutex                sync.RWMutex

	listenerWebsocket systemge.Listener[[]byte, systemge.Connection[[]byte]]
	httpServer        *httpServer.HTTPServer
	internalListener  systemge.Listener[*tools.Message, systemge.Connection[*tools.Message]]
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{
		requestResponseManager: tools.NewRequestResponseManager[*tools.Message](&configs.RequestResponseManager{}),
	}

	channelListener, err := listenerChannel.New[*tools.Message]("listenerChannel")
	if err != nil {
		panic(err)
	}
	app.internalListener = channelListener

	channelAccepter, err := serviceAccepter.New(
		channelListener,
		&configs.Accepter{},
		&configs.Routine{},
		func(connection systemge.Connection[*tools.Message]) error {

			_, err := serviceReader.NewAsync(
				connection,
				&configs.ReaderAsync{},
				&configs.Routine{},
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
				return err
			}

			app.mutex.Lock()
			app.internalConnections[connection] = struct{}{}
			app.mutex.Unlock()

			go func() { // abstract on close handler
				<-connection.GetCloseChannel()
				app.mutex.Lock()
				delete(app.internalConnections, connection)
				app.mutex.Unlock()
			}()

			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	app.channelAccepter = channelAccepter

	app.httpServer = httpServer.New(
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

	websocketAccepter, err := serviceAccepter.New(
		listenerWebsocket,
		&configs.Accepter{},
		&configs.Routine{},
		func(connection systemge.Connection[[]byte]) error {

			_, err := serviceTypedReader.NewSync(
				connection,
				&configs.ReaderSync{},
				&configs.Routine{},
				func(message *tools.Message, connection systemge.Connection[[]byte]) (*tools.Message, error) {
					app.mutex.RLock()
					defer app.mutex.RUnlock()

					for internalConnection := range app.internalConnections {
						go internalConnection.Write(message, 0)
					}
					return nil, nil
				},
				func(data []byte) (*tools.Message, error) {
					return tools.DeserializeMessage(data)
				},
				func(message *tools.Message) ([]byte, error) {
					return message.Serialize(), nil
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
				delete(app.websocketConnections, connection)
				app.mutex.Unlock()
			}()

			// propagate grid to new websocket connection

			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	app.websocketAccepter = websocketAccepter

	return app
}

/* func (app *AppWebsocketHTTP) OnConnectHandler(websocketClient *WebsocketServer.WebsocketClient) error {
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
*/
