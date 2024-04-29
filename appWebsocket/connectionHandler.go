package appWebsocket

import (
	"Systemge/WebsocketServer"
	"SystemgeSampleApp/typeDefinitions"
)

func (app *App) ConnectionHandler(connectionRequest *WebsocketServer.ConnectionRequest) {
	reponse, err := app.gridEndpoint.Request(typeDefinitions.GET_GRID_REQUEST.Create())
	if err != nil {
		return
	}
	connectionRequest.SendMessage(typeDefinitions.GET_GRID_WSPROPAGATE.Create([]string{reponse.Payload[0][0]}).ToBytes())
	app.websocketServer.AcceptConnectionRequest(connectionRequest, "")
}
