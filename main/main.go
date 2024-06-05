package main

import (
	"Systemge/Module"
	"SystemgeSampleApp/appGameOfLife"
	"SystemgeSampleApp/appWebsocket"
)

const TOPICRESOLUTIONSERVER_ADDRESS = ":60000"
const HTTP_DEV_PORT = ":8080"
const WEBSOCKET_PORT = ":8443"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	Module.StartSystemgeConsole(Module.NewMultiModule(
		Module.NewResolverServerFromConfig("topicResolutionServer.systemge", ERROR_LOG_FILE_PATH),
		Module.NewBrokerServerFromConfig("messageBrokerServerA.systemge", ERROR_LOG_FILE_PATH),
		Module.NewBrokerServerFromConfig("messageBrokerServerB.systemge", ERROR_LOG_FILE_PATH),
		Module.NewClient("clientGameOfLife", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, appGameOfLife.New),
		Module.NewWebsocketClient("clientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New),
		Module.NewHTTPServerFromConfig("httpServe.systemge", ERROR_LOG_FILE_PATH),
	), nil)
}
