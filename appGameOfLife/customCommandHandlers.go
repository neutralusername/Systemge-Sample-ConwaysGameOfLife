package appGameOfLife

import (
	"Systemge/Node"
	"Systemge/Utilities"
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topic"
)

func (app *App) GetCustomCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{
		"randomize":      app.randomizeGrid,
		"invert":         app.invertGrid,
		"chess":          app.chessGrid,
		"toggleToroidal": app.toggleToroidal,
	}
}

func (app *App) toggleToroidal(client *Node.Node, args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.toroidal = !app.toroidal
	if app.toroidal {
		println("Toroidal mode enabled")
	} else {
		println("Toroidal mode disabled")
	}
	return nil
}

func (app *App) randomizeGrid(client *Node.Node, args []string) error {
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
	err := client.AsyncMessage(topic.PROPGATE_GRID, client.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		client.GetLogger().Log(err.Error())
	}
	return nil
}

func (app *App) invertGrid(client *Node.Node, args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = 1 - app.grid[row][col]
		}
	}
	err := client.AsyncMessage(topic.PROPGATE_GRID, client.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		client.GetLogger().Log(err.Error())
	}
	return nil
}

func (app *App) chessGrid(client *Node.Node, args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = (row + col) % 2
		}
	}
	err := client.AsyncMessage(topic.PROPGATE_GRID, client.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		client.GetLogger().Log(err.Error())
	}
	return nil
}
