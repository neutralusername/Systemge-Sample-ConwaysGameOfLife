package main

import (
	"Systemge/Module"
	"SystemgeSampleApp/appGameOfLife"
	"SystemgeSampleApp/appWebsocket"
)

const TOPICRESOLUTIONSERVER_ADDRESS = ":60000"
const WEBSOCKET_PORT = ":8443"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Module.NewResolverServerFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH),
		Module.NewBrokerServerFromConfig("brokerGameOfLife.systemge", ERROR_LOG_FILE_PATH),
		Module.NewBrokerServerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH),
		Module.NewClient("clientGameOfLife", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, appGameOfLife.New),
		Module.NewWebsocketClient("clientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New),
		Module.NewHTTPServerFromConfig("httpServe.systemge", ERROR_LOG_FILE_PATH),
	), nil)
}
