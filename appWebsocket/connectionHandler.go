package appWebsocket

import (
	"Systemge/Websocket"
	"SystemgeSampleApp/typeDefinitions"
)

func (app *WebsocketApp) ConnectionHandler(connectionRequest *Websocket.ConnectionRequest) {
	connection := app.websocketServer.AcceptConnectionRequest(connectionRequest)
	err := app.messageBroker.Send(typeDefinitions.REQUEST_GRID_UNICAST.New([]string{connection.Id}))
	if err != nil {
		return
	}
}
