package appWebsocketHTTP

import (
	"SystemgeSampleConwaysGameOfLife/topics"

	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return app.websocketMessageHandlers
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	responseChannel, err := node.SyncMessage(topics.GET_GRID, websocketClient.GetId())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Failed to get grid", err).Error())
		}
		websocketClient.Disconnect()
		return
	}
	response, err := responseChannel.ReceiveResponse()
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Failed to receive response", err).Error())
		}
		websocketClient.Disconnect()
		return
	}
	getGridMessage := Message.NewAsync(topics.GET_GRID, response.GetPayload())
	websocketClient.Send([]byte(getGridMessage.Serialize()))
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {

}

func (app *AppWebsocketHTTP) propagateWebsocketAsyncMessage(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
	return node.AsyncMessage(message.GetTopic(), message.GetPayload())
}
