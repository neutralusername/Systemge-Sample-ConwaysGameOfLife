package appGameOfLife

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/Utilities"
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topic.GRID_CHANGE:     app.gridChange,
		topic.NEXT_GENERATION: app.nextGeneration,
		topic.SET_GRID:        app.setGrid,
	}
}

func (app *App) gridChange(client *Node.Node, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	gridChange := dto.UnmarshalGridChange(message.GetPayload())
	app.grid[gridChange.Row][gridChange.Column] = gridChange.State
	client.AsyncMessage(topic.PROPAGATE_GRID_CHANGE, client.GetName(), gridChange.Marshal())
	return nil
}

func (app *App) nextGeneration(client *Node.Node, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.calcNextGeneration()
	err := client.AsyncMessage(topic.PROPGATE_GRID, client.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		client.GetLogger().Log(Error.New("", err).Error())
	}
	return nil
}

func (app *App) setGrid(client *Node.Node, message *Message.Message) error {
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
	err := client.AsyncMessage(topic.PROPGATE_GRID, client.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		client.GetLogger().Log(Error.New("", err).Error())
	}
	return nil
}
