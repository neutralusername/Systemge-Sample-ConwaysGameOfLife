package appGameOfLife

import (
	"Systemge/Application"
	"Systemge/Message"
	"SystemgeSampleApp/dto"
	"SystemgeSampleApp/topic"
)

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{
		"randomize": app.randomizeGrid,
		"invert":    app.invertGrid,
	}
}

func (app *App) randomizeGrid(args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = app.randomizer.GenerateRandomNumber(0, 1)
		}
	}
	err := app.messageBrokerClient.AsyncMessage(Message.NewAsync(topic.GET_GRID, app.messageBrokerClient.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal()))
	if err != nil {
		app.logger.Log(err.Error())
	}
	return nil
}

func (app *App) invertGrid(args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = 1 - app.grid[row][col]
		}
	}
	err := app.messageBrokerClient.AsyncMessage(Message.NewAsync(topic.GET_GRID, app.messageBrokerClient.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal()))
	if err != nil {
		app.logger.Log(err.Error())
	}
	return nil
}
