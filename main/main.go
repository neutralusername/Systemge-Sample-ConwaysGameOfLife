package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Node"
	"Systemge/Resolver"
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocketHTTP"
	"SystemgeSampleConwaysGameOfLife/topics"
)

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	Node.StartCommandLineInterface(true,
		Node.New(Config.Node{
			Name: "nodeResolver",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Resolver.New(Config.Resolver{
			Server: Config.TcpServer{
				Port: 60000,
			},
			ConfigServer: Config.TcpServer{
				Port: 60001,
			},
			TcpTimeoutMs: 5000,
		})),
		Node.New(Config.Node{
			Name: "nodeBrokerGameOfLife",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Broker.New(Config.Broker{
			Server: Config.TcpServer{
				Port: 60002,
			},
			Endpoint: Config.TcpEndpoint{
				Address: "127.0.0.1:60002",
			},
			ConfigServer: Config.TcpServer{
				Port: 60003,
			},
			SyncTopics:  []string{topics.GET_GRID},
			AsyncTopics: []string{topics.GRID_CHANGE, topics.NEXT_GENERATION, topics.SET_GRID},
			ResolverConfigEndpoint: Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
			},
			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name: "nodeBrokerWebsocketHTTP",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Broker.New(Config.Broker{
			Server: Config.TcpServer{
				Port: 60004,
			},
			Endpoint: Config.TcpEndpoint{
				Address: "127.0.0.1:60004",
			},
			ConfigServer: Config.TcpServer{
				Port: 60005,
			},
			AsyncTopics: []string{topics.PROPAGATE_GRID_CHANGE, topics.PROPGATE_GRID},
			ResolverConfigEndpoint: Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
			},
			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name: "nodeGameOfLife",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, appGameOfLife.New()),
		Node.New(Config.Node{
			Name: "nodeWebsocketHTTP",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, appWebsocketHTTP.New()),
	)
}
