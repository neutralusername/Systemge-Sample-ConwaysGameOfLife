package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Utilities"
	"Systemge/Websocket"
)

func (app *App) OnConnectHandler(connection *Websocket.Connection) {
	response, err := app.messageBrokerClient.SyncMessage(Message.New("getGridSync", app.name, app.randomizer.GenerateRandomString(10, Utilities.ALPHA_NUMERIC), connection.Id))
	if err != nil {
		app.logger.Log(Error.New(err.Error()).Error())
		return
	}
	response.Topic = "getGrid"
	connection.Send([]byte(response.Serialize()))
}
