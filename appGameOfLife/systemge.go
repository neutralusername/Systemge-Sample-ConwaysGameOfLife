package appGameOfLife

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/Utilities"
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *App) GetSystemgeConfig() Config.Systemge {
	return Config.Systemge{
		HandleMessagesSequentially: false,
	}
}

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

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topic.GRID_CHANGE:     app.gridChange,
		topic.NEXT_GENERATION: app.nextGeneration,
		topic.SET_GRID:        app.setGrid,
	}
}

func (app *App) gridChange(node *Node.Node, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	gridChange := dto.UnmarshalGridChange(message.GetPayload())
	app.grid[gridChange.Row][gridChange.Column] = gridChange.State
	node.AsyncMessage(topic.PROPAGATE_GRID_CHANGE, node.GetName(), gridChange.Marshal())
	return nil
}

func (app *App) nextGeneration(node *Node.Node, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.calcNextGeneration()
	err := node.AsyncMessage(topic.PROPGATE_GRID, node.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		node.GetLogger().Error(Error.New("", err).Error())
	}
	return nil
}

func (app *App) setGrid(node *Node.Node, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if len(message.GetPayload()) != app.gridCols*app.gridRows {
		return Error.New("Invalid grid size", nil)
	}
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = Utilities.StringToInt(string(message.GetPayload()[row*app.gridCols+col]))
		}
	}
	err := node.AsyncMessage(topic.PROPGATE_GRID, node.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		node.GetLogger().Error(Error.New("", err).Error())
	}
	return nil
}
