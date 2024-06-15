package appWebsocket

import "Systemge/Application"

func (app *AppWebsocket) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{}
}
