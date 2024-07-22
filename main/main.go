package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Dashboard"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/Tools"
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocketHTTP"
	"SystemgeSampleConwaysGameOfLife/topics"
)

const LOGGER_PATH = "logs.log"

func main() {
	Node.New(&Config.Node{
		Name:           "dashboard",
		RandomizerSeed: Tools.GetSystemTime(),
	}, Dashboard.New(&Config.Dashboard{
		Server: &Config.TcpServer{
			Port:      8081,
			Blacklist: []string{},
			Whitelist: []string{},
		},
		StatusUpdateIntervalMs: 1000,
		HeapUpdateIntervalMs:   1000,
	},
		Node.New(&Config.Node{
			Name:           "nodeResolver",
			RandomizerSeed: Tools.GetSystemTime(),
		}, Resolver.New(&Config.Resolver{
			Server: &Config.TcpServer{
				Port: 60000,
			},
			ConfigServer: &Config.TcpServer{
				Port: 60001,
			},
			TcpTimeoutMs: 5000,
		})),
		Node.New(&Config.Node{
			Name:           "nodeBrokerGameOfLife",
			RandomizerSeed: Tools.GetSystemTime(),
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port: 60002,
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60002",
			},
			ConfigServer: &Config.TcpServer{
				Port: 60003,
			},
			SyncTopics:  []string{topics.GET_GRID},
			AsyncTopics: []string{topics.GRID_CHANGE, topics.NEXT_GENERATION, topics.SET_GRID},
			ResolverConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
			},
			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:           "nodeBrokerWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port: 60004,
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60004",
			},
			ConfigServer: &Config.TcpServer{
				Port: 60005,
			},
			AsyncTopics: []string{topics.PROPAGATE_GRID_CHANGE, topics.PROPGATE_GRID},
			ResolverConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
			},
			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:           "nodeGameOfLife",
			RandomizerSeed: Tools.GetSystemTime(),
		}, appGameOfLife.New()),
		Node.New(&Config.Node{
			Name:           "nodeWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
		}, appWebsocketHTTP.New()),
	)).Start()
	<-make(chan struct{})
}
