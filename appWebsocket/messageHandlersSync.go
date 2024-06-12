package appWebsocket

import "Systemge/Application"

func (app *WebsocketApp) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{}
}
