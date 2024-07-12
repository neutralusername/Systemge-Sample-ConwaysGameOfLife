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

const RESOLVER_ADDRESS = "127.0.0.1:60000"
const RESOLVER_NAME_INDICATION = "127.0.0.1"
const RESOLVER_TLS_CERT_PATH = "MyCertificate.crt"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

const ERROR_LOG_FILE_PATH = "error.log"

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
