package main

import (
	"Systemge/HTTP"
	"Systemge/MessageBrokerClient"
	"Systemge/MessageBrokerServer"
	"Systemge/Utilities"
	"Systemge/Websocket"
	"SystemgeSampleApp/appGameOfLife"
	"SystemgeSampleApp/appWebsocket"
	"time"
)

const MESSAGEBROKERSERVER_ADDRESS = ":60003"
const HTTP_DEV_PORT = ":8080"
const WEBSOCKET_PORT = ":8443"

func main() {
	logger := Utilities.NewLogger("error_log.txt")

	messageBrokerServer := MessageBrokerServer.New("messageBrokerServer", MESSAGEBROKERSERVER_ADDRESS, logger)
	messageBrokerServer.Start()

	messageBrokerClientWebsocket := MessageBrokerClient.New("messageBrokerClientWebsocket", logger, MESSAGEBROKERSERVER_ADDRESS)
	messageBrokerClientGameOfLife := MessageBrokerClient.New("messageBrokerClientGrid", logger, MESSAGEBROKERSERVER_ADDRESS)

	websocketServer := Websocket.New("websocketServer", messageBrokerClientWebsocket)
	appWebsocket := appWebsocket.New("websocketApp", logger, messageBrokerClientWebsocket, websocketServer)
	appGameOfLife := appGameOfLife.New("gameOfLifeApp", logger, messageBrokerClientGameOfLife)

	messageBrokerClientGameOfLife.Connect(appGameOfLife)
	messageBrokerClientWebsocket.Connect(appWebsocket)

	websocketServer.Start(appWebsocket)

	HTTPServerServe := HTTP.New(HTTP_DEV_PORT, "HTTPfrontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTP.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	HTTPServerWebsocket := HTTP.New(WEBSOCKET_PORT, "HTTPwebsocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTP.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	time.Sleep(1000000 * time.Second)
}
