package main

import (
	"Systemge/Application"
	"Systemge/HTTP"
	"Systemge/MessageBroker"
	"Systemge/MessageServerTCP"
	"Systemge/TCPServer"
	"Systemge/Utilities"
	"Systemge/Websocket"
	"SystemgeSampleApp/appGameOfLife"
	"SystemgeSampleApp/appWebsocket"
	"SystemgeSampleApp/typeDefinitions"
	"time"
)

const MESSAGEBROKER_ADDRESS = ":60003"

func main() {
	logger := Utilities.NewLogger("error_log.txt")

	tcpServerWebsocket := TCPServer.New(appWebsocket.ADDRESS, "tcpServerWebsocket")
	tcpServerWebsocket.Start()
	messageServerWebsocket := MessageServerTCP.New("messageServerWebsocket", tcpServerWebsocket, logger)
	/* messageServerWebsocket := MessageServerChannel.New("messageServerWebsocket", logger) */
	messageServerWebsocket.Start()

	tcpServerMessageBroker := TCPServer.New(MESSAGEBROKER_ADDRESS, "tcpServerMessageBroker")
	tcpServerMessageBroker.Start()
	messageServerMessageBroker := MessageServerTCP.New("messageServerMessageBroker", tcpServerMessageBroker, logger)
	/* messageServerMessageBroker := MessageServerChannel.New("messageServerMessageBroker", logger) */
	messageServerMessageBroker.Start()

	tcpServerGameOfLife := TCPServer.New(appGameOfLife.ADDRESS, "tcpServerGameOfLife")
	tcpServerGameOfLife.Start()
	messageServerGameOfLife := MessageServerTCP.New("messageServerGameOfLife", tcpServerGameOfLife, logger)
	/* messageServerGameOfLife := MessageServerChannel.New("messageServerGameOfLife", logger) */
	messageServerGameOfLife.Start()

	messageBroker := MessageBroker.New()
	subscriberWebsocket := MessageBroker.NewSubscriber(messageServerWebsocket.GetEndpoint(), true)
	subscriberGameOfLife := MessageBroker.NewSubscriber(messageServerGameOfLife.GetEndpoint(), true)
	messageBroker.AddSubscriber(subscriberWebsocket)
	messageBroker.AddSubscriber(subscriberGameOfLife)
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_BROADCAST, subscriberGameOfLife)
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_CHANGE, subscriberGameOfLife)
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_UNICAST, subscriberGameOfLife)
	messageBroker.AddMessageType(&typeDefinitions.WSPROPAGATE_MESSAGE_TYPE, subscriberWebsocket)

	messageBrokerServer := MessageBroker.NewServer("messageBrokerServer", messageBroker, messageServerMessageBroker, logger)
	messageBrokerServer.Start()

	websocketServer := Websocket.New("websocketServer", messageServerMessageBroker.GetEndpoint())

	appWebsocket := appWebsocket.New(websocketServer, messageServerMessageBroker.GetEndpoint())
	appServerWebsocket := Application.New("appServerWebsocket", logger, messageServerWebsocket)
	appServerWebsocket.Start(appWebsocket)

	appGameOfLife := appGameOfLife.New(messageServerMessageBroker.GetEndpoint(), logger)
	appServerGameOfLife := Application.New("appServerGameOfLife", logger, messageServerGameOfLife)
	appServerGameOfLife.Start(appGameOfLife)

	HTTPServerServe := HTTP.NewServer(HTTP.HTTP_DEV_PORT, "HTTPfrontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTP.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	HTTPServerWebsocket := HTTP.NewServer(HTTP.WEBSOCKET_PORT, "HTTPwebsocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTP.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	websocketServer.Start(appWebsocket)

	time.Sleep(1000000 * time.Second)
}
