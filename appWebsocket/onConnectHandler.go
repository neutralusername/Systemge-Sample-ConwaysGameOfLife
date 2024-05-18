package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Utilities"
	"Systemge/Websocket"
)

func (app *App) OnConnectHandler(connection *Websocket.Connection) {
	response, err := app.messageBrokerClient.SyncMessage(Message.New("getGridUnicast", app.name, app.randoizer.GenerateRandomString(10, Utilities.ALPHA_NUMERIC), connection.Id))
	if err != nil {
		app.logger.Log(Error.New("Failed to send getGrid message: " + err.Error()).Error())
	}
	connection.Send([]byte(response.Serialize()))
}
