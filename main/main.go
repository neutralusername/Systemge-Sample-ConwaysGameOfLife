package main

import (
	"Systemge/Application"
	"Systemge/HTTP"
	"Systemge/MessageBroker"
	"Systemge/MessageServerTCP"
	"Systemge/TCPServer"
	"Systemge/TypeDefinition"
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

	tcpServerWebsocket := TCPServer.New(appWebsocket.ADDRESS, "websocket")
	tcpServerWebsocket.Start()
	messageServerWebsocket := MessageServerTCP.New("websocket", tcpServerWebsocket, logger)
	/* messageServerWebsocket := MessageServerChannel.New("websocket", logger) */
	messageServerWebsocket.Start()

	tcpServerMessageBroker := TCPServer.New(MESSAGEBROKER_ADDRESS, "messageBroker")
	tcpServerMessageBroker.Start()
	messageServerMessageBroker := MessageServerTCP.New("messageBroker", tcpServerMessageBroker, logger)
	/* messageServerMessageBroker := MessageServerChannel.New("messageBroker", logger) */
	messageServerMessageBroker.Start()

	tcpServerGameOfLife := TCPServer.New(appGameOfLife.ADDRESS, "gameOfLife")
	tcpServerGameOfLife.Start()
	messageServerGameOfLife := MessageServerTCP.New("gameOfLife", tcpServerGameOfLife, logger)
	/* messageServerGameOfLife := MessageServerChannel.New("gameOfLife", logger) */
	messageServerGameOfLife.Start()

	messageBroker := MessageBroker.New()
	messageBroker.AddSubscriber(MessageBroker.NewSubscriber("websocket", messageServerWebsocket.GetEndpoint(), true))
	messageBroker.AddSubscriber(MessageBroker.NewSubscriber("gameOfLife", messageServerGameOfLife.GetEndpoint(), true))
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_BROADCAST, "gameOfLife")
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_CHANGE, "gameOfLife")
	messageBroker.AddMessageType(&typeDefinitions.REQUEST_GRID_UNICAST, "gameOfLife")
	messageBroker.AddMessageType(&TypeDefinition.WSPROPAGATE_MESSAGE_TYPE, "websocket")

	messageBrokerServer := MessageBroker.NewServer("messageBroker", messageBroker, messageServerMessageBroker, logger)
	messageBrokerServer.Start()

	websocketServer := Websocket.New("websocket", messageServerMessageBroker.GetEndpoint())

	appWebsocket := appWebsocket.New(websocketServer, messageServerMessageBroker.GetEndpoint())
	appServerWebsocket := Application.New("websocket", logger, messageServerWebsocket)
	appServerWebsocket.Start(appWebsocket)

	appGameOfLife := appGameOfLife.New(messageServerMessageBroker.GetEndpoint(), logger)
	appServerGameOfLife := Application.New("gameOfLife", logger, messageServerGameOfLife)
	appServerGameOfLife.Start(appGameOfLife)

	HTTPServerServe := HTTP.NewServer(HTTP.HTTP_DEV_PORT, "frontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTP.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	HTTPServerWebsocket := HTTP.NewServer(HTTP.WEBSOCKET_PORT, "websocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTP.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	websocketServer.Start(appWebsocket)

	time.Sleep(1000000 * time.Second)
}
