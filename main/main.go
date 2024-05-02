package main

import (
	"Systemge/ApplicationServer"
	"Systemge/HTTPServer"
	"Systemge/RequestServerTCP"
	"Systemge/TCPServer"
	"Systemge/Utilities"
	"Systemge/WebsocketServer"
	"SystemgeSampleApp/appGrid"
	"SystemgeSampleApp/appWebsocket"
	"time"
)

func main() {
	HTTPServerServe := HTTPServer.New(HTTPServer.HTTP_DEV_PORT, "frontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTPServer.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	logger := Utilities.NewLogger("error_log.txt")

	tcpServerWebsocket := TCPServer.New(appWebsocket.ADDRESS, "websocket")
	tcpServerWebsocket.Start()
	requestServerWebsocket := RequestServerTCP.New("websocket", tcpServerWebsocket, logger)
	requestServerWebsocket.Start()

	tcpServerGrid := TCPServer.New(appGrid.ADDRESS, "grid")
	tcpServerGrid.Start()
	requestServerGrid := RequestServerTCP.New("grid", tcpServerGrid, logger)
	requestServerGrid.Start()

	websocketServer := WebsocketServer.New("websocket")
	HTTPServerWebsocket := HTTPServer.New(HTTPServer.WEBSOCKET_PORT, "websocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTPServer.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()
	appServerWebsocket := ApplicationServer.New("websocket", logger, requestServerWebsocket)
	appWebsocket := appWebsocket.New(appServerWebsocket, websocketServer, requestServerGrid.GetEndpoint())
	websocketServer.Start(appWebsocket)

	appServerGrid := ApplicationServer.New("grid", logger, requestServerGrid)
	appGrid := appGrid.New(appServerGrid, requestServerWebsocket.GetEndpoint())

	appServerWebsocket.Start(appWebsocket)
	appServerGrid.Start(appGrid)

	time.Sleep(1000000 * time.Second)
}
