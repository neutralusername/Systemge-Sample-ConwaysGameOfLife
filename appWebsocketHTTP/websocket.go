package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleConwaysGameOfLife/topics"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return app.websocketMessageHandlers
}

func (app *AppWebsocketHTTP) GetWebsocketComponentConfig() *Config.Websocket {
	return app.websocketConfig
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	response, err := node.SyncMessage(topics.GET_GRID, node.GetName(), websocketClient.GetId())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Failed to get grid", err).Error())
		}
		websocketClient.Disconnect()
		return
	}
	websocketClient.Send([]byte(response.Serialize()))
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {

}

func (app *AppWebsocketHTTP) propagateWebsocketAsyncMessage(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
	return node.AsyncMessage(message.GetTopic(), message.GetOrigin(), message.GetPayload())
}
