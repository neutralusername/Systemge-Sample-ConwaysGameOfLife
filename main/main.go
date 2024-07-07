package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/TcpEndpoint"
	"Systemge/Utilities"
	"SystemgeSampleConwaysGameOfLife/appGameOfLife"
	"SystemgeSampleConwaysGameOfLife/appWebsocketHTTP"
	"SystemgeSampleConwaysGameOfLife/config"
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
			Name:                      config.NODE_GAMEOFLIFE_NAME,
			LoggerPath:                ERROR_LOG_FILE_PATH,
			ResolverEndpoint:          TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			SyncResponseTimeoutMs:     1000,
			TopicResolutionLifetimeMs: 10000,
			BrokerSubscribeDelayMs:    1000,
		}, appGameOfLife.New()),
		Node.New(Config.Node{
			Name:                      config.NODE_WEBSOCKET_HTTP_NAME,
			LoggerPath:                ERROR_LOG_FILE_PATH,
			ResolverEndpoint:          TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			SyncResponseTimeoutMs:     1000,
			TopicResolutionLifetimeMs: 10000,
			BrokerSubscribeDelayMs:    1000,
		}, applicationWebsocketHTTP),
	))
}
