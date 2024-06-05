package appWebsocket

import (
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"SystemgeSampleApp/topic"
)

func (app *App) GetWebsocketMessageHandlers() map[string]MessageBrokerClient.WebsocketMessageHandler {
	return map[string]MessageBrokerClient.WebsocketMessageHandler{
		topic.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
		topic.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
		topic.SET_GRID:        app.propagateWebsocketAsyncMessage,
	}
}
func (app *App) propagateWebsocketAsyncMessage(connection *MessageBrokerClient.WebsocketClient, message *Message.Message) error {
	return app.messageBrokerClient.AsyncMessage(message)
}

func (app *App) OnConnectHandler(connection *MessageBrokerClient.WebsocketClient) error {
	response, err := app.messageBrokerClient.SyncMessage(Message.NewSync(topic.GET_GRID_SYNC, app.messageBrokerClient.GetName(), connection.GetId()))
	if err != nil {
		return err
	}
	connection.Send([]byte(response.Serialize()))
	return nil
}

func (app *App) OnDisconnectHandler(connection *MessageBrokerClient.WebsocketClient) error {
	app.logger.Log("Connection closed")
	return nil
}
