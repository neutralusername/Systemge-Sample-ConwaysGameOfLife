package main

import (
	"Systemge/HTTP"
	"Systemge/Module"
	"SystemgeSampleApp/appGameOfLife"
	"SystemgeSampleApp/appWebsocket"
	"SystemgeSampleApp/topics"
	"net/http"
)

const MESSAGEBROKERSERVER_A_ADDRESS = ":60003"
const MESSAGEBROKERSERVER_B_ADDRESS = ":60004"
const TOPICRESOLUTIONSERVER_ADDRESS = ":60002"
const HTTP_DEV_PORT = ":8080"
const WEBSOCKET_PORT = ":8443"

const ERROR_LOG_FILE = "error_log.txt"

func main() {
	messageBrokerServerA := Module.NewMessageBrokerServerModule("messageBrokerServerA", MESSAGEBROKERSERVER_A_ADDRESS, ERROR_LOG_FILE,
		topics.GET_GRID_SYNC,
		topics.GRID_CHANGE,
		topics.NEXT_GENERATION,
		topics.SET_GRID,
	)
	messageBrokerServerB := Module.NewMessageBrokerServerModule("messageBrokerServerB", MESSAGEBROKERSERVER_B_ADDRESS, ERROR_LOG_FILE,
		topics.GET_GRID,
		topics.GET_GRID_CHANGE,
	)
	topicResolutionServer := Module.NewTopicResolutionModule("topicResolutionServer", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE, map[string]string{
		topics.GET_GRID_SYNC:   MESSAGEBROKERSERVER_A_ADDRESS,
		topics.GRID_CHANGE:     MESSAGEBROKERSERVER_A_ADDRESS,
		topics.NEXT_GENERATION: MESSAGEBROKERSERVER_A_ADDRESS,
		topics.SET_GRID:        MESSAGEBROKERSERVER_A_ADDRESS,
		topics.GET_GRID:        MESSAGEBROKERSERVER_B_ADDRESS,
		topics.GET_GRID_CHANGE: MESSAGEBROKERSERVER_B_ADDRESS,
	})
	httpServe := Module.NewHTTPModule("HTTPfrontend", HTTP_DEV_PORT, "", "", map[string]func(w http.ResponseWriter, r *http.Request){
		"/": HTTP.SendDirectory("../frontend"),
	})
	messageBrokerClientGameOfLife := Module.NewMessageBrokerClientModule("messageBrokerClientGameOfLife", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE, appGameOfLife.New)
	messageBrokerClientWebsocket := Module.NewWebsocketClientModule("messageBrokerClientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE, "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New)

	Module.CommandLoop(Module.NewMultiModule(
		messageBrokerServerA,
		messageBrokerServerB,
		topicResolutionServer,
		httpServe,
		messageBrokerClientGameOfLife,
		messageBrokerClientWebsocket,
	), messageBrokerClientGameOfLife.GetApplication().GetCustomCommandHandlers())
}
