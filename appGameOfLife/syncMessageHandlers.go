package appGameOfLife

import (
	"Systemge/Client"
	"Systemge/Message"
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *App) GetSyncMessageHandlers() map[string]Client.SyncMessageHandler {
	return map[string]Client.SyncMessageHandler{
		topic.GET_GRID: app.getGridSync,
	}
}

func (app *App) getGridSync(client *Client.Client, message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	return dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(), nil
}
