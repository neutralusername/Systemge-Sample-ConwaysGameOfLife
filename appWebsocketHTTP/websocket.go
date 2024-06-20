package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/Message"
	"Systemge/Utilities"
	"Systemge/WebsocketClient"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{
		topic.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
		topic.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
		topic.SET_GRID:        app.propagateWebsocketAsyncMessage,
	}
}
func (app *AppWebsocketHTTP) propagateWebsocketAsyncMessage(connection *WebsocketClient.Client, message *Message.Message) error {
	return app.client.AsyncMessage(message.GetTopic(), message.GetOrigin(), message.GetPayload())
}

func (app *AppWebsocketHTTP) OnConnectHandler(connection *WebsocketClient.Client) {
	response, err := app.client.SyncMessage(topic.GET_GRID, app.client.GetName(), connection.GetId())
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error sending sync message", err).Error())
	}
	connection.Send([]byte(response.Serialize()))
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(connection *WebsocketClient.Client) {
	app.client.GetLogger().Log("Connection closed")
}
