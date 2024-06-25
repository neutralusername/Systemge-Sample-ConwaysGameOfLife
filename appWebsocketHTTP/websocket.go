package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Client.WebsocketMessageHandler {
	return map[string]Client.WebsocketMessageHandler{
		topic.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
		topic.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
		topic.SET_GRID:        app.propagateWebsocketAsyncMessage,
	}
}
func (app *AppWebsocketHTTP) propagateWebsocketAsyncMessage(client *Client.Client, websocketClient *Client.WebsocketClient, message *Message.Message) error {
	return client.AsyncMessage(message.GetTopic(), message.GetOrigin(), message.GetPayload())
}

func (app *AppWebsocketHTTP) OnConnectHandler(client *Client.Client, websocketClient *Client.WebsocketClient) {
	response, err := client.SyncMessage(topic.GET_GRID, client.GetName(), websocketClient.GetId())
	if err != nil {
		client.GetLogger().Log(Utilities.NewError("Error sending sync message", err).Error())
	}
	websocketClient.Send([]byte(response.Serialize()))
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(client *Client.Client, websocketClient *Client.WebsocketClient) {
	client.GetLogger().Log("Connection closed")
}
