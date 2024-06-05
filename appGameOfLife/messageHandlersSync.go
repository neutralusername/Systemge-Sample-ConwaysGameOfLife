package appGameOfLife

import (
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"SystemgeSampleApp/dto"
	"SystemgeSampleApp/topic"
)

func (app *App) GetSyncMessageHandlers() map[string]MessageBrokerClient.SyncMessageHandler {
	return map[string]MessageBrokerClient.SyncMessageHandler{
		topic.GET_GRID_SYNC: app.getGridSync,
	}
}

func (app *App) getGridSync(message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	return dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(), nil
}
