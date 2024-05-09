package main

import (
	"Systemge/Application"
	"Systemge/HTTPServer"
	"Systemge/MessageBroker"
	"Systemge/RequestServerTCP"
	"Systemge/TCPServer"
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

	websocketServer := Websocket.New("websocket")
	HTTPServerWebsocket := HTTPServer.New(HTTPServer.WEBSOCKET_PORT, "websocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTPServer.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	tcpServerWebsocket := TCPServer.New(appWebsocket.ADDRESS, "websocket")
	tcpServerWebsocket.Start()
	requestServerWebsocket := RequestServerTCP.New("websocket", tcpServerWebsocket, logger)
	requestServerWebsocket.Start()

	tcpServerGrid := TCPServer.New(appGrid.ADDRESS, "grid")
	tcpServerGrid.Start()
	requestServerGrid := RequestServerTCP.New("grid", tcpServerGrid, logger)
	requestServerGrid.Start()

	tcpServerMessageBroker := TCPServer.New(MESSAGEBROKER_ADDRESS, "messageBroker")
	tcpServerMessageBroker.Start()
	requestServerMessageBroker := RequestServerTCP.New("messageBroker", tcpServerMessageBroker, logger)
	requestServerMessageBroker.Start()

	messageBroker := MessageBroker.New()
	messageBrokerServer := MessageBroker.NewServer("messageBroker", messageBroker, requestServerMessageBroker, logger)
	subscriberWebsocket := MessageBroker.NewSubscriber("websocket", requestServerWebsocket.GetEndpoint(), true)
	subscriberGrid := MessageBroker.NewSubscriber("grid", requestServerGrid.GetEndpoint(), true)
	messageBroker.AddSubscriber(subscriberWebsocket)
	messageBroker.AddSubscriber(subscriberGrid)
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_BROADCAST, "grid")
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_CHANGE, "grid")
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_UNICAST, "grid")
	messageBroker.AddMessageType(&typeDefinitions.BROADCAST_GRID, "websocket")
	messageBroker.AddMessageType(&typeDefinitions.BROADCAST_GRID_CHANGE, "websocket")
	messageBroker.AddMessageType(&typeDefinitions.UNICAST_GRID, "websocket")
	messageBrokerServer.Start()

	appWebsocket := appWebsocket.New(websocketServer, requestServerMessageBroker.GetEndpoint())
	appGrid := appGrid.New(requestServerMessageBroker.GetEndpoint(), logger)

	appServerWebsocket := Application.New("websocket", logger, requestServerWebsocket)
	appServerGrid := Application.New("grid", logger, requestServerGrid)

	websocketServer.Start(appWebsocket)
	appServerWebsocket.Start(appWebsocket)
	appServerGrid.Start(appGrid)

	time.Sleep(1000000 * time.Second)
}
