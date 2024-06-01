package appGameOfLife

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleApp/dto"
	"SystemgeSampleApp/topic"
)

func (app *App) GetGridSync(message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	return dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(), nil
}

func (app *App) GridChange(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	gridChange := dto.UnmarshalGridChange(message.Payload)
	app.grid[gridChange.Row][gridChange.Column] = gridChange.State
	app.messageBrokerClient.AsyncMessage(Message.NewAsync(topic.GET_GRID_CHANGE, app.messageBrokerClient.GetName(), gridChange.Marshal()))
	return nil
}

func (app *App) NextGeneration(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.calcNextGeneration()
	err := app.messageBrokerClient.AsyncMessage(Message.NewAsync(topic.GET_GRID, app.messageBrokerClient.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal()))
	if err != nil {
		app.logger.Log(Error.New("", err).Error())
	}
	return nil
}

func (app *App) SetGrid(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if len(message.Payload) != app.gridCols*app.gridRows {
		return Error.New("Invalid grid size", nil)
	}
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = Utilities.StringToInt(string(message.Payload[row*app.gridCols+col]))
		}
	}
	err := app.messageBrokerClient.AsyncMessage(Message.NewAsync(topic.GET_GRID, app.messageBrokerClient.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal()))
	if err != nil {
		app.logger.Log(Error.New("", err).Error())
	}
	return nil
}
