package appWebsocketHTTP

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		topic.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
		topic.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
		topic.SET_GRID:        app.propagateWebsocketAsyncMessage,
	}
}
func (app *AppWebsocketHTTP) propagateWebsocketAsyncMessage(client *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
	return client.AsyncMessage(message.GetTopic(), message.GetOrigin(), message.GetPayload())
}

func (app *AppWebsocketHTTP) OnConnectHandler(client *Node.Node, websocketClient *Node.WebsocketClient) {
	response, err := client.SyncMessage(topic.GET_GRID, client.GetName(), websocketClient.GetId())
	if err != nil {
		client.GetLogger().Log(Error.New("Error sending sync message", err).Error())
	}
	websocketClient.Send([]byte(response.Serialize()))
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(client *Node.Node, websocketClient *Node.WebsocketClient) {
	client.GetLogger().Log("Connection closed")
}
