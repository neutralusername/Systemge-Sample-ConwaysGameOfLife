package appWebsocket

import (
	"Systemge/WebsocketServer"
	"SystemgeSampleApp/typeDefinitions"
)

func (app *App) ConnectionHandler(connectionRequest *WebsocketServer.ConnectionRequest) {
	reponse, err := app.gridEndpoint.Request(typeDefinitions.GET_GRID_REQUEST.New())
	if err != nil {
		return
	}
	connectionRequest.SendMessage(typeDefinitions.GET_GRID_WSPROPAGATE.New([]string{reponse.Payload[0][0]}).Serialize())
	app.websocketServer.AcceptConnectionRequest(connectionRequest)
}
