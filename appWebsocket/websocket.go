package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/WebsocketClient"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *AppWebsocket) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{
		topic.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
		topic.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
		topic.SET_GRID:        app.propagateWebsocketAsyncMessage,
	}
}
func (app *AppWebsocket) propagateWebsocketAsyncMessage(connection *WebsocketClient.Client, message *Message.Message) error {
	return app.client.AsyncMessage(message.GetTopic(), message.GetOrigin(), message.GetPayload())
}

func (app *AppWebsocket) OnConnectHandler(connection *WebsocketClient.Client) {
	response, err := app.client.SyncMessage(topic.GET_GRID, app.client.GetName(), connection.GetId())
	if err != nil {
		app.client.GetLogger().Log(Error.New("Error sending sync message", err).Error())
	}
	connection.Send([]byte(response.Serialize()))
}

func (app *AppWebsocket) OnDisconnectHandler(connection *WebsocketClient.Client) {
	app.client.GetLogger().Log("Connection closed")
}
