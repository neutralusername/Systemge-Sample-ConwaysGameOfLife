package main

import (
	"Systemge/Application"
	"Systemge/HTTPServer"
	"Systemge/MessageBroker"
	"Systemge/MessageServerTCP"
	"Systemge/TCPServer"
	"Systemge/TypeDefinition"
	"Systemge/Utilities"
	"Systemge/Websocket"
	"SystemgeSampleApp/appGrid"
	"SystemgeSampleApp/appWebsocket"
	"SystemgeSampleApp/typeDefinitions"
	"time"
)

const MESSAGEBROKER_ADDRESS = ":60003"

func main() {
	logger := Utilities.NewLogger("error_log.txt")

	HTTPServerServe := HTTPServer.New(HTTPServer.HTTP_DEV_PORT, "frontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTPServer.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	tcpServerWebsocket := TCPServer.New(appWebsocket.ADDRESS, "websocket")
	tcpServerWebsocket.Start()
	messageServerWebsocket := MessageServerTCP.New("websocket", tcpServerWebsocket, logger)
	messageServerWebsocket.Start()

	tcpServerMessageBroker := TCPServer.New(MESSAGEBROKER_ADDRESS, "messageBroker")
	tcpServerMessageBroker.Start()
	messageServerMessageBroker := MessageServerTCP.New("messageBroker", tcpServerMessageBroker, logger)
	messageServerMessageBroker.Start()

	tcpServerGrid := TCPServer.New(appGrid.ADDRESS, "grid")
	tcpServerGrid.Start()
	messageServerGrid := MessageServerTCP.New("grid", tcpServerGrid, logger)
	messageServerGrid.Start()

	websocketServer := Websocket.New("websocket", messageServerWebsocket, messageServerMessageBroker.GetEndpoint())
	HTTPServerWebsocket := HTTPServer.New(HTTPServer.WEBSOCKET_PORT, "websocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTPServer.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	messageBroker := MessageBroker.New()
	messageBrokerServer := MessageBroker.NewServer("messageBroker", messageBroker, messageServerMessageBroker, logger)
	subscriberWebsocket := MessageBroker.NewSubscriber("websocket", messageServerWebsocket.GetEndpoint(), true)
	subscriberGrid := MessageBroker.NewSubscriber("grid", messageServerGrid.GetEndpoint(), true)
	messageBroker.AddSubscriber(subscriberWebsocket)
	messageBroker.AddSubscriber(subscriberGrid)
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_BROADCAST, "grid")
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_CHANGE, "grid")
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_UNICAST, "grid")
	messageBroker.AddMessageType(&TypeDefinition.WSPROPAGATE_MESSAGE_TYPE, "websocket")
	messageBrokerServer.Start()

	appWebsocket := appWebsocket.New(websocketServer, messageServerMessageBroker.GetEndpoint())
	appGrid := appGrid.New(messageServerMessageBroker.GetEndpoint(), logger)

	appServerGrid := Application.New("grid", logger, messageServerGrid)
	appServerWebsocket := Application.New("websocket", logger, messageServerWebsocket)

	websocketServer.Start(appWebsocket)
	appServerGrid.Start(appGrid)
	appServerWebsocket.Start(appWebsocket)

	time.Sleep(1000000 * time.Second)
}
