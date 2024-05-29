package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Websocket"
)

func (app *App) OnConnectHandler(connection *Websocket.Connection) {
	response, err := app.messageBrokerClient.SyncMessage(Message.NewSync("getGridSync", app.name, connection.Id))
	if err != nil {
		app.logger.Log(Error.New(err.Error()).Error())
		return
	}
	connection.Send([]byte(response.Serialize()))
}
