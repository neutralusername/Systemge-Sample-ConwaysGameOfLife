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
	applicationWebsocketHTTP := appWebsocketHTTP.New()
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Resolver.New(Config.ParseResolverConfigFromFile("resolver.systemge")),
		Broker.New(Config.ParseBrokerConfigFromFile("brokerGameOfLife.systemge")),
		Broker.New(Config.ParseBrokerConfigFromFile("brokerWebsocketHTTP.systemge")),
		Node.New(Config.Node{
			Name:       "nodeGameOfLife",
			LoggerPath: ERROR_LOG_FILE_PATH,
		}, appGameOfLife.New(), nil, nil),
		Node.New(Config.Node{
			Name:       "nodeWebsocketHTTP",
			LoggerPath: ERROR_LOG_FILE_PATH,
		}, applicationWebsocketHTTP, applicationWebsocketHTTP, applicationWebsocketHTTP),
	))
}
