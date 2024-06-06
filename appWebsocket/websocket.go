package appWebsocket

import (
	"Systemge/Application"
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

func (app *App) OnConnectHandler(connection *WebsocketClient.Client) error {
	response, err := app.messageBrokerClient.SyncMessage(Message.NewSync(topic.GET_GRID_SYNC, app.messageBrokerClient.GetName(), connection.GetId()))
	if err != nil {
		return err
	}
	connection.Send([]byte(response.Serialize()))
	return nil
}

func (app *App) OnDisconnectHandler(connection *WebsocketClient.Client) error {
	app.logger.Log("Connection closed")
	return nil
}
