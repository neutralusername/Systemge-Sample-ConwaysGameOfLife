package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Client"
)

type WebsocketApp struct {
	client *Client.Client
}

func New(client *Client.Client, args []string) Application.WebsocketApplication {
	return &WebsocketApp{
		client: client,
	}
}

func (app *WebsocketApp) OnStart() error {
	return nil
}

func (app *WebsocketApp) OnStop() error {
	return nil
}
