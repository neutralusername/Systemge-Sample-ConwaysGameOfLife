package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/WebsocketClient"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *WebsocketApp) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{
		topic.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
		topic.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
		topic.SET_GRID:        app.propagateWebsocketAsyncMessage,
	}
}
func (app *WebsocketApp) propagateWebsocketAsyncMessage(connection *WebsocketClient.Client, message *Message.Message) error {
	return app.messageBrokerClient.AsyncMessage(message.GetTopic(), message.GetOrigin(), message.GetPayload())
}

func (app *WebsocketApp) OnConnectHandler(connection *WebsocketClient.Client) {
	response, err := app.messageBrokerClient.SyncMessage(topic.GET_GRID, app.messageBrokerClient.GetName(), connection.GetId())
	if err != nil {
		app.logger.Log(Error.New("Error sending sync message", err).Error())
	}
	connection.Send([]byte(response.Serialize()))
}

func (app *WebsocketApp) OnDisconnectHandler(connection *WebsocketClient.Client) {
	app.logger.Log("Connection closed")
}
