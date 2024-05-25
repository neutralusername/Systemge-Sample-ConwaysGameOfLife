package main

import (
	"Systemge/HTTP"
	"Systemge/MessageBrokerClient"
	"Systemge/MessageBrokerServer"
	"Systemge/TopicResolutionServer"
	"Systemge/Utilities"
	"Systemge/Websocket"
	"SystemgeSampleApp/appGameOfLife"
	"SystemgeSampleApp/appWebsocket"
	"time"
)

const MESSAGEBROKERSERVER_ADDRESS = ":60003"
const TOPICRESOLUTIONSERVER_ADDRESS = ":60002"
const HTTP_DEV_PORT = ":8080"
const WEBSOCKET_PORT = ":8443"

func main() {
	logger := Utilities.NewLogger("error_log.txt")

	topicResolutionServer := TopicResolutionServer.New("topicResolutionServer", TOPICRESOLUTIONSERVER_ADDRESS, logger)
	topicResolutionServer.Start()
	topicResolutionServer.RegisterTopic("getGridSync", MESSAGEBROKERSERVER_ADDRESS)
	topicResolutionServer.RegisterTopic("gridChange", MESSAGEBROKERSERVER_ADDRESS)
	topicResolutionServer.RegisterTopic("nextGeneration", MESSAGEBROKERSERVER_ADDRESS)
	topicResolutionServer.RegisterTopic("setGrid", MESSAGEBROKERSERVER_ADDRESS)
	topicResolutionServer.RegisterTopic("getGrid", MESSAGEBROKERSERVER_ADDRESS)
	topicResolutionServer.RegisterTopic("getGridChange", MESSAGEBROKERSERVER_ADDRESS)

	messageBrokerServer := MessageBrokerServer.New("messageBrokerServer", MESSAGEBROKERSERVER_ADDRESS, logger)
	messageBrokerServer.AddTopic("getGridSync")
	messageBrokerServer.AddTopic("gridChange")
	messageBrokerServer.AddTopic("nextGeneration")
	messageBrokerServer.AddTopic("setGrid")
	messageBrokerServer.AddTopic("getGrid")
	messageBrokerServer.AddTopic("getGridChange")
	messageBrokerServer.Start()

	messageBrokerClientWebsocket := MessageBrokerClient.New("messageBrokerClientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, logger)
	messageBrokerClientWebsocket.Connect(MESSAGEBROKERSERVER_ADDRESS)

	messageBrokerClientGameOfLife := MessageBrokerClient.New("messageBrokerClientGrid", TOPICRESOLUTIONSERVER_ADDRESS, logger)
	messageBrokerClientGameOfLife.Connect(MESSAGEBROKERSERVER_ADDRESS)

	websocketServer := Websocket.New("websocketServer", messageBrokerClientWebsocket)
	appWebsocket := appWebsocket.New("websocketApp", logger, messageBrokerClientWebsocket, websocketServer)
	appGameOfLife := appGameOfLife.New("gameOfLifeApp", logger, messageBrokerClientGameOfLife, 90, 140)

	messageBrokerClientGameOfLife.RegisterIncomingSyncTopic("getGridSync", appGameOfLife.GetGridSync)
	messageBrokerClientGameOfLife.RegisterIncomingAsyncTopic("gridChange", appGameOfLife.GridChange)
	messageBrokerClientGameOfLife.RegisterIncomingAsyncTopic("nextGeneration", appGameOfLife.NextGeneration)
	messageBrokerClientGameOfLife.RegisterIncomingAsyncTopic("setGrid", appGameOfLife.SetGrid)

	messageBrokerClientWebsocket.RegisterIncomingAsyncTopic("getGrid", appWebsocket.GetGrid)
	messageBrokerClientWebsocket.RegisterIncomingAsyncTopic("getGridChange", appWebsocket.GetGridChange)

	websocketServer.SetOnMessageHandler(appWebsocket.OnMessageHandler)
	websocketServer.SetOnConnectHandler(appWebsocket.OnConnectHandler)

	websocketServer.Start()

	HTTPServerServe := HTTP.New(HTTP_DEV_PORT, "HTTPfrontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTP.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	HTTPServerWebsocket := HTTP.New(WEBSOCKET_PORT, "HTTPwebsocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTP.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	time.Sleep(1000000 * time.Second)
}
