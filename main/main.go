package main

import (
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocketHttp"
	"time"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/DashboardServer"
)

const LOGGER_PATH = "logs.log"

func main() {
	dashboardServer := DashboardServer.New("dashboardServer",
		&Config.DashboardServer{
			HTTPServerConfig: &Config.HTTPServer{
				TcpServerConfig: &Config.TcpServer{
					Port: 8081,
				},
			},
			WebsocketServerConfig: &Config.WebsocketServer{
				Pattern:                 "/ws",
				ClientWatchdogTimeoutMs: 1000 * 60,
				TcpServerConfig: &Config.TcpServer{
					Port: 8444,
				},
			},
			SystemgeServerConfig: &Config.SystemgeServer{
				TcpSystemgeListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						Port: 60000,
					},
				},
				TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
			},
			DashboardSystemgeCommands:   true,
			DashboardHttpCommands:       true,
			DashboardWebsocketCommands:  true,
			FrontendHeartbeatIntervalMs: 1000 * 60,
			UpdateIntervalMs:            1000,
			MaxMetricEntries:            100,
		},
		nil, nil,
	)
	if err := dashboardServer.Start(); err != nil {
		panic(err)
	}
	appWebsocketHttp.New()
	appGameOfLife.New()
	<-make(chan time.Time)
}
