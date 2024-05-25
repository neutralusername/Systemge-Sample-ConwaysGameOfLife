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

const MESSAGEBROKERSERVER_A_ADDRESS = ":60003"
const MESSAGEBROKERSERVER_B_ADDRESS = ":60004"
const TOPICRESOLUTIONSERVER_ADDRESS = ":60002"
const HTTP_DEV_PORT = ":8080"
const WEBSOCKET_PORT = ":8443"

const GET_GRID_SYNC_TOPIC = "getGridSync"
const GRID_CHANGE_TOPIC = "gridChange"
const NEXT_GENERATION_TOPIC = "nextGeneration"
const SET_GRID_TOPIC = "setGrid"
const GET_GRID_TOPIC = "getGrid"
const GET_GRID_CHANGE_TOPIC = "getGridChange"

func main() {
	logger := Utilities.NewLogger("error_log.txt")

	messageBrokerServerA := MessageBrokerServer.New("messageBrokerServerA", MESSAGEBROKERSERVER_A_ADDRESS, logger)
	messageBrokerServerA.Start()
	messageBrokerServerA.AddTopics(GET_GRID_SYNC_TOPIC, GRID_CHANGE_TOPIC, NEXT_GENERATION_TOPIC, SET_GRID_TOPIC)

	messageBrokerServerB := MessageBrokerServer.New("messageBrokerServerB", MESSAGEBROKERSERVER_B_ADDRESS, logger)
	messageBrokerServerB.Start()
	messageBrokerServerB.AddTopics(GET_GRID_TOPIC, GET_GRID_CHANGE_TOPIC)

	topicResolutionServer := TopicResolutionServer.New("topicResolutionServer", TOPICRESOLUTIONSERVER_ADDRESS, logger)
	topicResolutionServer.Start()
	topicResolutionServer.RegisterTopics(MESSAGEBROKERSERVER_A_ADDRESS, GET_GRID_SYNC_TOPIC, GRID_CHANGE_TOPIC, NEXT_GENERATION_TOPIC, SET_GRID_TOPIC)
	topicResolutionServer.RegisterTopics(MESSAGEBROKERSERVER_B_ADDRESS, GET_GRID_TOPIC, GET_GRID_CHANGE_TOPIC)

	messageBrokerClientWebsocket := MessageBrokerClient.New("messageBrokerClientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, logger)
	messageBrokerClientWebsocket.Connect()

	messageBrokerClientGameOfLife := MessageBrokerClient.New("messageBrokerClientGrid", TOPICRESOLUTIONSERVER_ADDRESS, logger)
	messageBrokerClientGameOfLife.Connect()

	websocketServer := Websocket.New("websocketServer", messageBrokerClientWebsocket)
	websocketServer.Start()

	appWebsocket := appWebsocket.New("websocketApp", logger, messageBrokerClientWebsocket, websocketServer)
	appGameOfLife := appGameOfLife.New("gameOfLifeApp", logger, messageBrokerClientGameOfLife, 90, 140)

	messageBrokerClientGameOfLife.RegisterIncomingSyncTopic(GET_GRID_SYNC_TOPIC, appGameOfLife.GetGridSync)
	messageBrokerClientGameOfLife.RegisterIncomingAsyncTopic(GRID_CHANGE_TOPIC, appGameOfLife.GridChange)
	messageBrokerClientGameOfLife.RegisterIncomingAsyncTopic(NEXT_GENERATION_TOPIC, appGameOfLife.NextGeneration)
	messageBrokerClientGameOfLife.RegisterIncomingAsyncTopic(SET_GRID_TOPIC, appGameOfLife.SetGrid)

	messageBrokerClientWebsocket.RegisterIncomingAsyncTopic(GET_GRID_TOPIC, appWebsocket.GetGrid)
	messageBrokerClientWebsocket.RegisterIncomingAsyncTopic(GET_GRID_CHANGE_TOPIC, appWebsocket.GetGridChange)

	websocketServer.SetOnMessageHandler(appWebsocket.OnMessageHandler)
	websocketServer.SetOnConnectHandler(appWebsocket.OnConnectHandler)

	HTTPServerServe := HTTP.New(HTTP_DEV_PORT, "HTTPfrontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTP.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	HTTPServerWebsocket := HTTP.New(WEBSOCKET_PORT, "HTTPwebsocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTP.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	time.Sleep(1000000 * time.Second)
}
