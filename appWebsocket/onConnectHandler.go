package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Websocket"
)

func (app *App) OnConnectHandler(connection *Websocket.Connection) {
	err := app.messageBrokerClient.Send(Message.New("getGridUnicast", app.name, connection.Id))
	if err != nil {
		app.logger.Log(Error.New("Failed to send getGridUnicast message: " + err.Error()).Error())
	}
}
