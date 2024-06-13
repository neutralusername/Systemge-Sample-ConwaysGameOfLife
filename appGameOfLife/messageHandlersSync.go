package appGameOfLife

import (
	"Systemge/Application"
	"Systemge/Message"
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topic.GET_GRID_SYNC: app.getGridSync,
	}
}

func (app *App) getGridSync(message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	return dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(), nil
}
