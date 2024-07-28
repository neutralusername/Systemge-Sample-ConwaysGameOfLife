package appBrokerWebsocketHttp

import (
	"SystemgeSampleConwaysGameOfLife/topics"

	"github.com/neutralusername/Systemge/Config"
)

func (app *App) GetBrokerComponentConfig() *Config.Broker {
	return &Config.Broker{
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
	}
}
