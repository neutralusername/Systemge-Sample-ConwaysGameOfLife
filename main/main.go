package main

import (
	"Systemge/HTTP"
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"Systemge/MessageBrokerServer"
	"Systemge/Utilities"
	"Systemge/Websocket"
	"SystemgeSampleApp/appGameOfLife"
	"SystemgeSampleApp/appWebsocket"
	"time"
)

const MESSAGEBROKERSERVER_ADDRESS = ":60003"

func main() {
	logger := Utilities.NewLogger("error_log.txt")

	messageBrokerServer := MessageBrokerServer.New("messageBrokerServer", MESSAGEBROKERSERVER_ADDRESS, logger)
	messageBrokerServer.Start()

	messageBrokerClientWebsocket := MessageBrokerClient.New("messageBrokerClientWebsocket")
	messageBrokerClientGameOfLife := MessageBrokerClient.New("messageBrokerClientGrid")

	websocketServer := Websocket.New("websocketServer")

	appWebsocket := appWebsocket.New("websocketApp", logger, messageBrokerClientWebsocket, websocketServer)
	appGameOfLife := appGameOfLife.New("gameOfLifeApp", logger, messageBrokerClientGameOfLife)

	messageHandlersGameOfLife := map[string]func(*Message.Message) error{
		"gridChange":     appGameOfLife.GridChange,
		"getGridUnicast": appGameOfLife.GetGridUnicast,
	}
	messageHandlersWebsocket := map[string]func(*Message.Message) error{
		"websocketUnicast": appWebsocket.WebsocketUnicast,
		"getGrid":          appWebsocket.GetGrid,
		"getGridChange":    appWebsocket.GetGridChange,
	}

	messageBrokerClientGameOfLife.Connect(MESSAGEBROKERSERVER_ADDRESS, messageHandlersGameOfLife)
	messageBrokerClientWebsocket.Connect(MESSAGEBROKERSERVER_ADDRESS, messageHandlersWebsocket)
	websocketServer.Start(appWebsocket)

	HTTPServerServe := HTTP.NewServer(HTTP.HTTP_DEV_PORT, "HTTPfrontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTP.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	HTTPServerWebsocket := HTTP.NewServer(HTTP.WEBSOCKET_PORT, "HTTPwebsocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTP.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	time.Sleep(1000000 * time.Second)
}
