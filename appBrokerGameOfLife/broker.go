package appBrokerGameOfLife

import (
	"SystemgeSampleConwaysGameOfLife/topics"

	"github.com/neutralusername/Systemge/Config"
)

func (app *App) GetBrokerComponentConfig() *Config.Broker {
	return &Config.Broker{
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
	}
}
