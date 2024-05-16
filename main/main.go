package main

import (
	"Systemge/Application"
	"Systemge/HTTP"
	"Systemge/MessageBrokerClient"
	"Systemge/MessageBrokerServer"
	"Systemge/Utilities"
	"Systemge/Websocket"
	"SystemgeSampleApp/appGameOfLife"
	WebsocketApp "SystemgeSampleApp/websocketApp"
	"time"
)

const MESSAGEBROKER_ADDRESS = ":60003"

func main() {
	logger := Utilities.NewLogger("error_log.txt")

	messageBrokerServer := MessageBrokerServer.New("messageBrokerServer", MESSAGEBROKER_ADDRESS, logger)
	messageBrokerServer.Start()

	messageBrokerServer.AddMessageType("gridChange")
	messageBrokerServer.AddMessageType("getGridUnicast")
	messageBrokerServer.AddMessageType("websocketUnicast")
	messageBrokerServer.AddMessageType("getGridChange")
	messageBrokerServer.AddMessageType("getGrid")

	messageBrokerClientGameOfLife := MessageBrokerClient.New("messageBrokerClientGrid")
	messageBrokerClientGameOfLife.Connect(MESSAGEBROKER_ADDRESS)
	messageBrokerClientGameOfLife.Subscribe("gridChange")
	messageBrokerClientGameOfLife.Subscribe("getGridUnicast")

	messageBrokerClientWebsocket := MessageBrokerClient.New("messageBrokerClientWebsocket")
	messageBrokerClientWebsocket.Connect(MESSAGEBROKER_ADDRESS)
	messageBrokerClientWebsocket.Subscribe("websocketUnicast")
	messageBrokerClientWebsocket.Subscribe("getGrid")
	messageBrokerClientWebsocket.Subscribe("getGridChange")

	websocketServer := Websocket.New("websocketServer", messageBrokerClientWebsocket)
	websocketApp := WebsocketApp.New("websocketApp", logger, messageBrokerClientWebsocket, websocketServer)
	websocketServer.Start(websocketApp)

	applicationServerGameOfLife := Application.New("applicationServerGameOfLife", logger, messageBrokerClientGameOfLife)
	applicationServerGameOfLife.Start(appGameOfLife.New("gameOfLifeApp", logger, messageBrokerClientGameOfLife))

	applicationServerWebsocket := Application.New("applicationServerWebsocket", logger, messageBrokerClientWebsocket)
	applicationServerWebsocket.Start(websocketApp)

	HTTPServerServe := HTTP.NewServer(HTTP.HTTP_DEV_PORT, "HTTPfrontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTP.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	HTTPServerWebsocket := HTTP.NewServer(HTTP.WEBSOCKET_PORT, "HTTPwebsocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTP.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	time.Sleep(1000000 * time.Second)
}
