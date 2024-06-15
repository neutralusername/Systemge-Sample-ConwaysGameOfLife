package main

import (
	"Systemge/Module"
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocket"
)

const TOPICRESOLUTIONSERVER_ADDRESS = "127.0.0.1:60000"
const WEBSOCKET_PORT = ":8443"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	clientGameOfLife := Module.NewClient("clientGameOfLife", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, appGameOfLife.New, nil)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Module.NewResolverFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH),
		Module.NewBrokerFromConfig("brokerGameOfLife.systemge", ERROR_LOG_FILE_PATH),
		Module.NewBrokerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH),
		clientGameOfLife,
		Module.NewWebsocketClient("clientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New, nil),
		Module.NewHTTPServerFromConfig("httpServe.systemge", ERROR_LOG_FILE_PATH),
	), clientGameOfLife.GetApplication().GetCustomCommandHandlers())
}
