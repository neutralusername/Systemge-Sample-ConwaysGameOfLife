package appWebsocket

import (
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"SystemgeSampleApp/topics"
)

func (app *App) GetOnMessageHandler() MessageBrokerClient.OnMessageHandler {
	return func(connection *MessageBrokerClient.WebsocketClient, message *Message.Message) {
		message.Origin = connection.Id
		message.SyncRequestToken = ""
		message.SyncResponseToken = ""
		switch message.Topic {
		case topics.GRID_CHANGE:
			err := app.messageBrokerClient.AsyncMessage(message)
			if err != nil {
				connection.Send([]byte(Message.NewAsync("error", app.messageBrokerClient.GetName(), err.Error()).Serialize()))
			}
		case topics.NEXT_GENERATION:
			err := app.messageBrokerClient.AsyncMessage(message)
			if err != nil {
				connection.Send([]byte(Message.NewAsync("error", app.messageBrokerClient.GetName(), err.Error()).Serialize()))
			}
		case topics.SET_GRID:
			err := app.messageBrokerClient.AsyncMessage(message)
			if err != nil {
				connection.Send([]byte(Message.NewAsync("error", app.messageBrokerClient.GetName(), err.Error()).Serialize()))
			}
		default:
			connection.Send([]byte(Message.NewAsync("error", app.messageBrokerClient.GetName(), "Unknown message topic \""+message.Topic+"\"").Serialize()))
		}
	}
}
