package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/TcpServer"
	"SystemgeSampleConwaysGameOfLife/topics"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		topics.GRID_CHANGE:     app.propagateWebsocketAsyncMessage,
		topics.NEXT_GENERATION: app.propagateWebsocketAsyncMessage,
		topics.SET_GRID:        app.propagateWebsocketAsyncMessage,
	}
}
func (app *AppWebsocketHTTP) propagateWebsocketAsyncMessage(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
	return node.AsyncMessage(message.GetTopic(), message.GetOrigin(), message.GetPayload())
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	response, err := node.SyncMessage(topics.GET_GRID, node.GetName(), websocketClient.GetId())
	if err != nil {
		node.GetLogger().Error(Error.New("Error sending sync message", err).Error())
		websocketClient.Disconnect()
		return
	}
	websocketClient.Send([]byte(response.Serialize()))
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	node.GetLogger().Info("Client \"" + websocketClient.GetId() + "\" disconnected")
}

func (app *AppWebsocketHTTP) GetWebsocketComponentConfig() Config.Websocket {
	return Config.Websocket{
		Pattern:                          "/ws",
		Server:                           TcpServer.New(8443, "", ""),
		HandleClientMessagesSequentially: false,
		ClientMessageCooldownMs:          0,
		ClientWatchdogTimeoutMs:          20000,
	}
}
