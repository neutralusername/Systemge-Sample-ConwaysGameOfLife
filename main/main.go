package main

import (
	"Systemge/Module"
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
	clientGameOfLife := Module.NewClient(&Module.ClientConfig{
		Name:                   "clientGameOfLife",
		ResolverAddress:        RESOLVER_ADDRESS,
		ResolverNameIndication: RESOLVER_NAME_INDICATION,
		ResolverTLSCertPath:    RESOLVER_TLS_CERT_PATH,
		LoggerPath:             ERROR_LOG_FILE_PATH,
	}, appGameOfLife.New, nil)
	clientWebsocketHTTP := Module.NewCompositeClientWebsocketHTTP(&Module.ClientConfig{
		Name:                   "clientWebsocketHTTP",
		ResolverAddress:        RESOLVER_ADDRESS,
		ResolverNameIndication: RESOLVER_NAME_INDICATION,
		ResolverTLSCertPath:    RESOLVER_TLS_CERT_PATH,
		LoggerPath:             ERROR_LOG_FILE_PATH,
		WebsocketPattern:       "/ws",
		WebsocketPort:          WEBSOCKET_PORT,
		HTTPPort:               HTTP_PORT,
	}, appWebsocketHTTP.New, nil)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Module.NewResolverFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH),
		Module.NewBrokerFromConfig("brokerGameOfLife.systemge", ERROR_LOG_FILE_PATH),
		Module.NewBrokerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH),
		clientGameOfLife,
		clientWebsocketHTTP,
	), clientGameOfLife.GetApplication().GetCustomCommandHandlers(), clientWebsocketHTTP.GetCustomCommandHandlers())
}
