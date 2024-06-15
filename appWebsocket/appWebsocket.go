package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Client"
)

type AppWebsocket struct {
	client *Client.Client
}

func New(client *Client.Client, args []string) (Application.WebsocketApplication, error) {
	return &AppWebsocket{
		client: client,
	}, nil
}

func (app *AppWebsocket) OnStart() error {
	return nil
}

func (app *AppWebsocket) OnStop() error {
	return nil
}
