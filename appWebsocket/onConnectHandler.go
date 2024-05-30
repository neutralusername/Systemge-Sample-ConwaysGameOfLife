package appWebsocket

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"SystemgeSampleApp/topics"
)

func (app *App) GetOnConnectHandler() MessageBrokerClient.OnConnectHandler {
	return func(connection *MessageBrokerClient.WebsocketClient) {
		response, err := app.messageBrokerClient.SyncMessage(Message.NewSync(topics.GET_GRID_SYNC, app.messageBrokerClient.GetName(), connection.Id))
		if err != nil {
			app.logger.Log(Error.New(err.Error()).Error())
			return
		}
		connection.Send([]byte(response.Serialize()))
	}
}
