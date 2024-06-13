package appGameOfLife

import (
	"Systemge/Application"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topic.GRID_CHANGE:     app.gridChange,
		topic.NEXT_GENERATION: app.nextGeneration,
		topic.SET_GRID:        app.setGrid,
	}
}

func (app *App) gridChange(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	gridChange := dto.UnmarshalGridChange(message.GetPayload())
	app.grid[gridChange.Row][gridChange.Column] = gridChange.State
	app.messageBrokerClient.AsyncMessage(Message.NewAsync(topic.GET_GRID_CHANGE, app.messageBrokerClient.GetName(), gridChange.Marshal()))
	return nil
}

func (app *App) nextGeneration(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.calcNextGeneration()
	err := app.messageBrokerClient.AsyncMessage(Message.NewAsync(topic.GET_GRID, app.messageBrokerClient.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal()))
	if err != nil {
		app.logger.Log(Error.New("", err).Error())
	}
	return nil
}

func (app *App) setGrid(message *Message.Message) error {
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
	err := app.messageBrokerClient.AsyncMessage(Message.NewAsync(topic.GET_GRID, app.messageBrokerClient.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal()))
	if err != nil {
		app.logger.Log(Error.New("", err).Error())
	}
	return nil
}
