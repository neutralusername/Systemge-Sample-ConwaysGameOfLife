package appWebsocket

import (
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
)

func (app *App) GetOnMessageHandler() MessageBrokerClient.OnMessageHandler {
	return func(connection *MessageBrokerClient.WebsocketClient, message *Message.Message) {
		message.Origin = connection.Id
		message.SyncRequestToken = ""
		message.SyncResponseToken = ""
		switch message.Topic {
		case "heartbeat":
			connection.ResetWatchdog()
		case "gridChange":
			err := app.messageBrokerClient.AsyncMessage(message)
			if err != nil {
				connection.Send([]byte(Message.NewAsync("error", app.messageBrokerClient.GetName(), err.Error()).Serialize()))
			}
		case "nextGeneration":
			err := app.messageBrokerClient.AsyncMessage(message)
			if err != nil {
				connection.Send([]byte(Message.NewAsync("error", app.messageBrokerClient.GetName(), err.Error()).Serialize()))
			}
		case "setGrid":
			err := app.messageBrokerClient.AsyncMessage(message)
			if err != nil {
				connection.Send([]byte(Message.NewAsync("error", app.messageBrokerClient.GetName(), err.Error()).Serialize()))
			}
		default:
			connection.Send([]byte(Message.NewAsync("error", app.messageBrokerClient.GetName(), "Unknown message type").Serialize()))
		}
	}
}
