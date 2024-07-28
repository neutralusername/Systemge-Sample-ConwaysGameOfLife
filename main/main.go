package main

import (
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocketHTTP"
	"SystemgeSampleConwaysGameOfLife/topics"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Tools"
)

const LOGGER_PATH = "logs.log"

func main() {
	loggerQueue := Tools.NewLoggerQueue(LOGGER_PATH, 10000)
	dashboardNode := Node.New(&Config.Node{
		Name:           "dashboard",
		RandomizerSeed: Tools.GetSystemTime(),
	}, Dashboard.New(&Config.Dashboard{
		Server: &Config.TcpServer{
			Port:      8081,
			Blacklist: []string{},
			Whitelist: []string{},
		},
		NodeStatusIntervalMs:           1000,
		NodeSystemgeCounterIntervalMs:  1000,
		NodeWebsocketCounterIntervalMs: 1000,
		NodeBrokerCounterIntervalMs:    1000,
		NodeResolverCounterIntervalMs:  1000,
		HeapUpdateIntervalMs:           1000,
		GoroutineUpdateIntervalMs:      1000,
		NodeSpawnerCounterIntervalMs:   1000,
		NodeHTTPCounterIntervalMs:      1000,
		AutoStart:                      true,
		AddDashboardToDashboard:        true,
	},
		Node.New(&Config.Node{
			Name:           "nodeResolver",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeResolver\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeResolver\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeResolver\"] ", loggerQueue),
		}, Node.NewResolverApplication(&Config.Resolver{
			Server: &Config.TcpServer{
				Port: 60000,
			},
			ConfigServer: &Config.TcpServer{
				Port: 60001,
			},
			TcpTimeoutMs: 5000,
		})),
		Node.New(&Config.Node{
			Name:           "nodeBrokerWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeBrokerWebsocketHTTP\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeBrokerWebsocketHTTP\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeBrokerWebsocketHTTP\"] ", loggerQueue),
		}, Node.NewBrokerApplication(&Config.Broker{
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
			Name:           "nodeBrokerGameOfLife",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeBrokerGameOfLife\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeBrokerGameOfLife\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeBrokerGameOfLife\"] ", loggerQueue),
		}, Node.NewBrokerApplication(&Config.Broker{
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
			Name:           "nodeWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeWebsocketHTTP\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeWebsocketHTTP\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeWebsocketHTTP\"] ", loggerQueue),
		}, appWebsocketHTTP.New()),
		Node.New(&Config.Node{
			Name:           "nodeGameOfLife",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeGameOfLife\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeGameOfLife\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeGameOfLife\"] ", loggerQueue),
		}, appGameOfLife.New()),
	))
	dashboardNode.StartBlocking()
}
