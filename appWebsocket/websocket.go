package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/WebsocketClient"
	"SystemgeSampleApp/topic"
)

func (app *App) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{
		topic.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
		topic.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
		topic.SET_GRID:        app.propagateWebsocketAsyncMessage,
	}
}
func (app *App) propagateWebsocketAsyncMessage(connection *WebsocketClient.Client, message *Message.Message) error {
	return app.messageBrokerClient.AsyncMessage(message)
}

func (app *App) OnConnectHandler(connection *WebsocketClient.Client) {
	response, err := app.messageBrokerClient.SyncMessage(Message.NewSync(topic.GET_GRID_SYNC, app.messageBrokerClient.GetName(), connection.GetId()))
	if err != nil {
		app.logger.Log(Error.New("Error sending sync message", err).Error())
	}
	connection.Send([]byte(response.Serialize()))
}

func (app *App) OnDisconnectHandler(connection *WebsocketClient.Client) {
	app.logger.Log("Connection closed")
}
