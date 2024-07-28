package appResolver

import "github.com/neutralusername/Systemge/Config"

func (app *App) GetResolverComponentConfig() *Config.Resolver {
	return &Config.Resolver{
		Server: &Config.TcpServer{
			Port: 60000,
		},
		ConfigServer: &Config.TcpServer{
			Port: 60001,
		},
		TcpTimeoutMs: 5000,
	}
}
