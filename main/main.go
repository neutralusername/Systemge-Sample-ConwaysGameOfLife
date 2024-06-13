package main

import (
	"Systemge/Module"
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocket"
)

const TOPICRESOLUTIONSERVER_ADDRESS = ":60000"
const WEBSOCKET_PORT = ":8443"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	clientGameOfLife := Module.NewClient("clientGameOfLife", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, appGameOfLife.New)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Module.NewResolverServerFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH),
		Module.NewBrokerServerFromConfig("brokerGameOfLife.systemge", ERROR_LOG_FILE_PATH),
		Module.NewBrokerServerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH),
		clientGameOfLife,
		Module.NewWebsocketClient("clientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New),
		Module.NewHTTPServerFromConfig("httpServe.systemge", ERROR_LOG_FILE_PATH),
	), clientGameOfLife.GetApplication().GetCustomCommandHandlers())
}
