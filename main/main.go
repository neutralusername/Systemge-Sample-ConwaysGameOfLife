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
			MaxChartEntries:           100,
			HeapUpdateIntervalMs:      1000,
			GoroutineUpdateIntervalMs: 1000,
			StatusUpdateIntervalMs:    1000,
			MetricsUpdateIntervalMs:   1000,
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
