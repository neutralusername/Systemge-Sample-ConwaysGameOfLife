package main

import (
	"Systemge/HTTP"
	"Systemge/Module"
	"SystemgeSampleApp/appGameOfLife"
	"SystemgeSampleApp/appWebsocket"
	"SystemgeSampleApp/topic"
	"net/http"
)

const MESSAGEBROKERSERVER_A_ADDRESS = ":60003"
const MESSAGEBROKERSERVER_B_ADDRESS = ":60004"
const TOPICRESOLUTIONSERVER_ADDRESS = ":60002"
const HTTP_DEV_PORT = ":8080"
const WEBSOCKET_PORT = ":8443"

const ERROR_LOG_FILE = "error_log.txt"

func main() {
	messageBrokerServerA := Module.NewMessageBrokerServer("messageBrokerServerA", MESSAGEBROKERSERVER_A_ADDRESS, ERROR_LOG_FILE,
		topic.GET_GRID_SYNC,
		topic.GRID_CHANGE,
		topic.NEXT_GENERATION,
		topic.SET_GRID,
	)
	messageBrokerServerB := Module.NewMessageBrokerServer("messageBrokerServerB", MESSAGEBROKERSERVER_B_ADDRESS, ERROR_LOG_FILE,
		topic.GET_GRID,
		topic.GET_GRID_CHANGE,
	)
	topicResolutionServer := Module.NewTopicResolutionServer("topicResolutionServer", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE, map[string]string{
		topic.GET_GRID_SYNC:   MESSAGEBROKERSERVER_A_ADDRESS,
		topic.GRID_CHANGE:     MESSAGEBROKERSERVER_A_ADDRESS,
		topic.NEXT_GENERATION: MESSAGEBROKERSERVER_A_ADDRESS,
		topic.SET_GRID:        MESSAGEBROKERSERVER_A_ADDRESS,
		topic.GET_GRID:        MESSAGEBROKERSERVER_B_ADDRESS,
		topic.GET_GRID_CHANGE: MESSAGEBROKERSERVER_B_ADDRESS,
	})
	httpServe := Module.NewHTTPServer("HTTPfrontend", HTTP_DEV_PORT, "", "", map[string]func(w http.ResponseWriter, r *http.Request){
		"/": HTTP.SendDirectory("../frontend"),
	})
	messageBrokerClientGameOfLife := Module.NewMessageBrokerClient("messageBrokerClientGameOfLife", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE, appGameOfLife.New)
	messageBrokerClientWebsocket := Module.NewWebsocketClient("messageBrokerClientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE, "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New)

	Module.CommandLoop(Module.NewMultiModule(
		messageBrokerServerA,
		messageBrokerServerB,
		topicResolutionServer,
		httpServe,
		messageBrokerClientGameOfLife,
		messageBrokerClientWebsocket,
	), messageBrokerClientGameOfLife.GetApplication().GetCustomCommandHandlers())
}
