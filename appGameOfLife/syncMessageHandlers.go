package appGameOfLife

import (
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topic.GET_GRID: app.getGridSync,
	}
}

func (app *App) getGridSync(node *Node.Node, message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	return dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(), nil
}
