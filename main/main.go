package main

import (
	"Systemge/HTTP"
	"Systemge/Module"
	"SystemgeSampleApp/appGameOfLife"
	"SystemgeSampleApp/appWebsocket"
	"net/http"
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
	messageBrokerServerA := Module.NewServerModule("messageBrokerServerA", MESSAGEBROKERSERVER_A_ADDRESS, "error_log.txt",
		GET_GRID_SYNC_TOPIC, GRID_CHANGE_TOPIC, NEXT_GENERATION_TOPIC, SET_GRID_TOPIC,
	)
	messageBrokerServerB := Module.NewServerModule("messageBrokerServerB", MESSAGEBROKERSERVER_B_ADDRESS, "error_log.txt",
		GET_GRID_TOPIC, GET_GRID_CHANGE_TOPIC,
	)
	topicResolutionServer := Module.NewResolutionModule("topicResolutionServer", TOPICRESOLUTIONSERVER_ADDRESS, "error_log.txt", map[string]string{
		GET_GRID_SYNC_TOPIC:   MESSAGEBROKERSERVER_A_ADDRESS,
		GRID_CHANGE_TOPIC:     MESSAGEBROKERSERVER_A_ADDRESS,
		NEXT_GENERATION_TOPIC: MESSAGEBROKERSERVER_A_ADDRESS,
		SET_GRID_TOPIC:        MESSAGEBROKERSERVER_A_ADDRESS,
		GET_GRID_TOPIC:        MESSAGEBROKERSERVER_B_ADDRESS,
		GET_GRID_CHANGE_TOPIC: MESSAGEBROKERSERVER_B_ADDRESS,
	})
	httpServe := Module.NewHTTPModule("HTTPfrontend", HTTP_DEV_PORT, "", "", map[string]func(w http.ResponseWriter, r *http.Request){
		"/": HTTP.SendDirectory("../frontend"),
	})
	messageBrokerClientGameOfLife := Module.NewClientModule("messageBrokerClientGameOfLife", TOPICRESOLUTIONSERVER_ADDRESS, "error_log.txt", appGameOfLife.New)
	messageBrokerClientWebsocket := Module.NewWebsocketClientModule("messageBrokerClientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, "error_log.txt", "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New)

	Module.CommandLoop(Module.NewMultiModule(
		messageBrokerServerA,
		messageBrokerServerB,
		topicResolutionServer,
		httpServe,
		messageBrokerClientGameOfLife,
		messageBrokerClientWebsocket,
	), nil)
}
