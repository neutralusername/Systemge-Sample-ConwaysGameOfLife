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
	HTTPServerServe := HTTPServer.Create(HTTPServer.HTTP_DEV_PORT, "frontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTPServer.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	logger := Utilities.CreateLogger("error_log.txt")

	tcpServerWebsocket := TCPServer.Create(appWebsocket.ADDRESS, "websocket")
	tcpServerWebsocket.Start()
	requestServerWebsocket := RequestServerTCP.Create("websocket", tcpServerWebsocket, logger)
	requestServerWebsocket.Start()

	tcpServerGrid := TCPServer.Create(appGrid.ADDRESS, "grid")
	tcpServerGrid.Start()
	requestServerGrid := RequestServerTCP.Create("grid", tcpServerGrid, logger)
	requestServerGrid.Start()

	websocketServer := WebsocketServer.Create("websocket")
	HTTPServerWebsocket := HTTPServer.Create(HTTPServer.WEBSOCKET_PORT, "websocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTPServer.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()
	appServerWebsocket := ApplicationServer.Create("websocket", logger, requestServerWebsocket)
	appWebsocket := appWebsocket.Create(appServerWebsocket, websocketServer, requestServerGrid.GetEndpoint())
	websocketServer.Start(appWebsocket)

	appServerGrid := ApplicationServer.Create("grid", logger, requestServerGrid)
	appGrid := appGrid.Create(appServerGrid, requestServerWebsocket.GetEndpoint())

	appServerWebsocket.Start(appWebsocket)
	appServerGrid.Start(appGrid)

	time.Sleep(1000000 * time.Second)
}
