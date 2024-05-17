package main

import (
	"Systemge/HTTP"
	"Systemge/Message"
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

	messageBrokerClientWebsocket := MessageBrokerClient.New("messageBrokerClientWebsocket")
	messageBrokerClientGameOfLife := MessageBrokerClient.New("messageBrokerClientGrid")

	websocketServer := Websocket.New("websocketServer")
	websocketApp := WebsocketApp.New("websocketApp", logger, messageBrokerClientWebsocket, websocketServer)

	gameOfLifeApp := appGameOfLife.New("gameOfLifeApp", logger, messageBrokerClientGameOfLife)

	messageHandlersGameOfLife := map[string]func(*Message.Message) error{
		"gridChange":     gameOfLifeApp.GridChange,
		"getGridUnicast": gameOfLifeApp.GetGridUnicast,
	}

	messageHandlersWebsocket := map[string]func(*Message.Message) error{
		"websocketUnicast": websocketApp.WebsocketUnicast,
		"getGrid":          websocketApp.GetGrid,
		"getGridChange":    websocketApp.GetGridChange,
	}

	messageBrokerClientGameOfLife.Connect(MESSAGEBROKER_ADDRESS, messageHandlersGameOfLife)
	messageBrokerClientWebsocket.Connect(MESSAGEBROKER_ADDRESS, messageHandlersWebsocket)

	websocketServer.Start(websocketApp)

	HTTPServerServe := HTTP.NewServer(HTTP.HTTP_DEV_PORT, "HTTPfrontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTP.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	HTTPServerWebsocket := HTTP.NewServer(HTTP.WEBSOCKET_PORT, "HTTPwebsocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTP.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	time.Sleep(1000000 * time.Second)
}
