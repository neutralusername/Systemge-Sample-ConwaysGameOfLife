package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/TcpEndpoint"
	"Systemge/TcpServer"
	"Systemge/Utilities"
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocketHTTP"
	"SystemgeSampleConwaysGameOfLife/topics"
)

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	Module.StartCommandLineInterface(Module.NewMultiModule(true,
		Node.New(Config.Node{
			Name:   "nodeResolver",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Resolver.New(Config.Resolver{
			Server:       TcpServer.New(60000, "", ""),
			ConfigServer: TcpServer.New(60001, "", ""),
			TcpTimeoutMs: 5000,
		})),
		Node.New(Config.Node{
			Name:   "nodeBrokerGameOfLife",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Broker.New(Config.Broker{
			Server:       TcpServer.New(60002, "", ""),
			Endpoint:     TcpEndpoint.New("127.0.0.1:60002", "", ""),
			ConfigServer: TcpServer.New(60003, "", ""),

			SyncTopics:  []string{topics.GET_GRID},
			AsyncTopics: []string{topics.GRID_CHANGE, topics.NEXT_GENERATION, topics.SET_GRID},

			ResolverConfigEndpoint: TcpEndpoint.New("127.0.0.1:60001", "", ""),

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name:   "nodeBrokerWebsocketHTTP",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Broker.New(Config.Broker{
			Server:       TcpServer.New(60004, "", ""),
			Endpoint:     TcpEndpoint.New("127.0.0.1:60004", "", ""),
			ConfigServer: TcpServer.New(60005, "", ""),

			AsyncTopics: []string{topics.PROPAGATE_GRID_CHANGE, topics.PROPGATE_GRID},

			ResolverConfigEndpoint: TcpEndpoint.New("127.0.0.1:60001", "", ""),

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name:   "nodeGameOfLife",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, appGameOfLife.New()),
		Node.New(Config.Node{
			Name:   "nodeWebsocketHTTP",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, appWebsocketHTTP.New()),
	))
}
