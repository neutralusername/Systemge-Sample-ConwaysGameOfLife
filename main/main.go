package main

import (
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocketHttp"
	"time"
)

const LOGGER_PATH = "logs.log"

func main() {
	app := appGameOfLife.NewApp(appGameOfLife.NewChannel())
	appWebsocketHttp.New(app.Listener.GetConnector())
	<-make(chan time.Time)
}

/* dashboardServer := DashboardServer.New("dashboardServer",
	&Config.DashboardServer{
		HTTPServerConfig: &Config.HTTPServer{
			TcpServerConfig: &Config.TcpServer{
				Port: 8081,
			},
		},
		WebsocketServerConfig: &Config.WebsocketServer{
			Pattern:                 "/ws",
			ClientWatchdogTimeoutMs: 1000 * 60 * 3,
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
		FrontendHeartbeatIntervalMs: 1000 * 60 * 1,
		UpdateIntervalMs:            1000,
		MaxEntriesPerMetrics:        100,
	},
	nil, nil,
)
if err := dashboardServer.Start(); err != nil {
	panic(err)
} */
