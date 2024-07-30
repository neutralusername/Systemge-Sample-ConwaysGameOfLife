package main

import (
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocketHTTP"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Tools"
)

const LOGGER_PATH = "logs.log"

func main() {
	Tools.NewLoggerQueue(LOGGER_PATH, 10000)
	dashboardNode := Node.New(&Config.Node{
		Name:           "dashboard",
		RandomizerSeed: Tools.GetSystemTime(),
	}, Dashboard.New(&Config.Dashboard{
		ServerConfig: &Config.TcpServer{
			Port:      8081,
			Blacklist: []string{},
			Whitelist: []string{},
		},
		NodeStatusIntervalMs:           1000,
		NodeSystemgeCounterIntervalMs:  1000,
		NodeWebsocketCounterIntervalMs: 1000,
		HeapUpdateIntervalMs:           1000,
		GoroutineUpdateIntervalMs:      1000,
		NodeSpawnerCounterIntervalMs:   1000,
		NodeHTTPCounterIntervalMs:      1000,
		AutoStart:                      true,
		AddDashboardToDashboard:        true,
	},
		Node.New(&Config.Node{
			Name:              "nodeWebsocketHTTP",
			RandomizerSeed:    Tools.GetSystemTime(),
			InfoLoggerPath:    LOGGER_PATH,
			WarningLoggerPath: LOGGER_PATH,
			ErrorLoggerPath:   LOGGER_PATH,
		}, appWebsocketHTTP.New()),
		Node.New(&Config.Node{
			Name:              "nodeGameOfLife",
			RandomizerSeed:    Tools.GetSystemTime(),
			InfoLoggerPath:    LOGGER_PATH,
			WarningLoggerPath: LOGGER_PATH,
			ErrorLoggerPath:   LOGGER_PATH,
		}, appGameOfLife.New()),
	))
	dashboardNode.StartBlocking()
}
