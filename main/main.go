package main

import (
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocketHTTP"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Tools"
)

const LOGGER_PATH = "logs.log"

func main() {
	Tools.NewLoggerQueue(LOGGER_PATH, 10000)
	dashboardNode := Dashboard.New(&Config.Dashboard{
		NodeConfig: &Config.Node{
			Name:           "dashboard",
			RandomizerSeed: Tools.GetSystemTime(),
		},
		ServerConfig: &Config.TcpServer{
			Port:      8081,
			Blacklist: []string{},
			Whitelist: []string{},
		},
		NodeStatusIntervalMs:                                    1000,
		NodeSystemgeCounterIntervalMs:                           1000,
		NodeSystemgeInvalidMessageCounterIntervalMs:             1000,
		NodeSystemgeIncomingConnectionAttemptsCounterIntervalMs: 1000,
		NodeSystemgeIncomingSyncResponseCounterIntervalMs:       1000,
		NodeSystemgeIncomingSyncRequestCounterIntervalMs:        1000,
		NodeSystemgeIncomingAsyncMessageCounterIntervalMs:       1000,
		NodeSystemgeOutgoingConnectionAttemptCounterIntervalMs:  1000,
		NodeSystemgeOutgoingSyncRequestCounterIntervalMs:        1000,
		NodeSystemgeOutgoingAsyncMessageCounterIntervalMs:       1000,
		NodeSystemgeOutgoingSyncResponsesCounterIntervalMs:      1000,
		NodeWebsocketCounterIntervalMs:                          1000,
		HeapUpdateIntervalMs:                                    1000,
		GoroutineUpdateIntervalMs:                               1000,
		NodeSpawnerCounterIntervalMs:                            1000,
		NodeHTTPCounterIntervalMs:                               1000,
		AutoStart:                                               true,
		AddDashboardToDashboard:                                 true,
	},
		Node.New(&Config.NewNode{
			NodeConfig: &Config.Node{Name: "nodeWebsocketHTTP",
				RandomizerSeed:    Tools.GetSystemTime(),
				InfoLoggerPath:    LOGGER_PATH,
				WarningLoggerPath: LOGGER_PATH,
				ErrorLoggerPath:   LOGGER_PATH,
			},
			SystemgeConfig: &Config.Systemge{
				ProcessMessagesOfEachConnectionSequentially: true,
				ProcessAllMessagesSequentially:              false,

				SyncRequestTimeoutMs:            10000,
				TcpTimeoutMs:                    5000,
				MaxConnectionAttempts:           0,
				ConnectionAttemptDelayMs:        1000,
				StopAfterOutgoingConnectionLoss: true,
				ServerConfig: &Config.TcpServer{
					Port: 60002,
				},
				EndpointConfigs: []*Config.TcpEndpoint{
					{
						Address: "localhost:60001",
					},
				},
				IncomingMessageByteLimit: 0,
				MaxPayloadSize:           0,
				MaxTopicSize:             0,
				MaxSyncTokenSize:         0,
			},
			HttpConfig: &Config.HTTP{
				ServerConfig: &Config.TcpServer{
					Port: 8080,
				},
			},
			WebsocketConfig: &Config.Websocket{
				Pattern: "/ws",
				ServerConfig: &Config.TcpServer{
					Port:      8443,
					Blacklist: []string{},
					Whitelist: []string{},
				},
				HandleClientMessagesSequentially: false,
				ClientWatchdogTimeoutMs:          20000,
				Upgrader: &websocket.Upgrader{
					ReadBufferSize:  1024,
					WriteBufferSize: 1024,
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
			},
		}, appWebsocketHTTP.New()),
		Node.New(&Config.NewNode{
			NodeConfig: &Config.Node{
				Name:              "nodeGameOfLife",
				RandomizerSeed:    Tools.GetSystemTime(),
				InfoLoggerPath:    LOGGER_PATH,
				WarningLoggerPath: LOGGER_PATH,
				ErrorLoggerPath:   LOGGER_PATH,
			},
			SystemgeConfig: appGameOfLife.SYSTEMGE_CONFIG,
		}, appGameOfLife.New()),
	)
	dashboardNode.StartBlocking()
}
