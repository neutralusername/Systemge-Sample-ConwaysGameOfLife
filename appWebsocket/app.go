package appWebsocket

import (
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"SystemgeSampleApp/topic"
)

type App struct {
	logger              *Utilities.Logger
	messageBrokerClient *MessageBrokerClient.Client
}

func New(logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client) MessageBrokerClient.WebsocketApplication {
	return &App{
		logger:              logger,
		messageBrokerClient: messageBrokerClient,
	}
}

func (app *App) GetAsyncMessageHandlers() map[string]MessageBrokerClient.AsyncMessageHandler {
	return map[string]MessageBrokerClient.AsyncMessageHandler{
		topic.GET_GRID:        app.WebsocketPropagate,
		topic.GET_GRID_CHANGE: app.WebsocketPropagate,
	}
}
func (app *App) WebsocketPropagate(message *Message.Message) error {
	app.messageBrokerClient.WebsocketBroadcast([]byte(message.Serialize()))
	return nil
}

func (app *App) GetSyncMessageHandlers() map[string]MessageBrokerClient.SyncMessageHandler {
	return map[string]MessageBrokerClient.SyncMessageHandler{}
}

func (app *App) GetCustomCommandHandlers() map[string]func([]string) error {
	return map[string]func([]string) error{}
}

func (app *App) GetWebsocketMessageHandlers() map[string]MessageBrokerClient.WebsocketMessageHandler {
	return map[string]MessageBrokerClient.WebsocketMessageHandler{
		topic.GRID_CHANGE:     app.PropagateWebsocketAsyncMessage,
		topic.NEXT_GENERATION: app.PropagateWebsocketAsyncMessage,
		topic.SET_GRID:        app.PropagateWebsocketAsyncMessage,
	}
}
func (app *App) PropagateWebsocketAsyncMessage(connection *MessageBrokerClient.WebsocketClient, message *Message.Message) error {
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
