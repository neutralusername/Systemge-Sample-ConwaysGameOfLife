package appWebsocketHTTP

import (
	"net/http"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/HTTP"
)

func (app *AppWebsocketHTTP) GetHTTPMessageHandlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/": HTTP.SendDirectory("../frontend"),
	}
}

func (app *AppWebsocketHTTP) GetHTTPComponentConfig() *Config.HTTP {
	return &Config.HTTP{
		ServerConfig: &Config.TcpServer{
			Port: 8080,
		},
	}
}
