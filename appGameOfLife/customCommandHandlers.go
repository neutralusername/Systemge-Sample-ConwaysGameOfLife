package appGameOfLife

import (
	"Systemge/Application"
	"Systemge/Utilities"
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{
		"randomize": app.randomizeGrid,
		"invert":    app.invertGrid,
		"chess":     app.chessGrid,
	}
}

func (app *App) randomizeGrid(args []string) error {
	percentageOfAliveCells := 50
	if len(args) > 0 {
		percentageOfAliveCells = Utilities.StringToInt(args[0])
	}
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			if app.randomizer.GenerateRandomNumber(1, 100) <= percentageOfAliveCells {
				app.grid[row][col] = 1
			} else {
				app.grid[row][col] = 0
			}
		}
	}
	err := app.messageBrokerClient.AsyncMessage(topic.PROPGATE_GRID, app.messageBrokerClient.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
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
	err := app.messageBrokerClient.AsyncMessage(topic.PROPGATE_GRID, app.messageBrokerClient.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		app.logger.Log(err.Error())
	}
	return nil
}

func (app *App) chessGrid(args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = (row + col) % 2
		}
	}
	err := app.messageBrokerClient.AsyncMessage(topic.PROPGATE_GRID, app.messageBrokerClient.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		app.logger.Log(err.Error())
	}
	return nil
}
