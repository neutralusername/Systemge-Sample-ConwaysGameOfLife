package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Resolver"
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocketHTTP"
)

func main() {
	err := Resolver.New(Config.ParseResolverConfigFromFile("resolver.systemge")).Start()
	if err != nil {
		panic(err)
	}
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Broker.New(Config.ParseBrokerConfigFromFile("brokerGameOfLife.systemge")),
		Broker.New(Config.ParseBrokerConfigFromFile("brokerWebsocketHTTP.systemge")),
		Node.New(Config.ParseNodeConfigFromFile("nodeGameOfLife.systemge"), appGameOfLife.New()),
		Node.New(Config.ParseNodeConfigFromFile("nodeWebsocketHTTP.systemge"), appWebsocketHTTP.New()),
	))
}
