package appGameOfLife

import (
	"Systemge/Helpers"
	"Systemge/Node"
	"Systemge/Tools"
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"
)

type App struct {
	randomizer *Tools.Randomizer
	grid       [][]int
	mutex      sync.Mutex
	gridRows   int
	gridCols   int
	toroidal   bool
}

func New() *App {
	app := &App{
		randomizer: Tools.NewRandomizer(Tools.GetSystemTime()),

		grid:     nil,
		gridRows: 90,
		gridCols: 140,
		toroidal: true,
	}
	return app
}

func (app *App) GetCommandHandlers() map[string]Node.CommandHandler {
	return map[string]Node.CommandHandler{
		"randomize":      app.randomizeGrid,
		"invert":         app.invertGrid,
		"chess":          app.chessGrid,
		"toggleToroidal": app.toggleToroidal,
	}
}

func (app *App) toggleToroidal(node *Node.Node, args []string) error {
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

func (app *App) randomizeGrid(node *Node.Node, args []string) error {
	percentageOfAliveCells := 50
	if len(args) > 0 {
		percentageOfAliveCells = Helpers.StringToInt(args[0])
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
	err := node.AsyncMessage(topics.PROPGATE_GRID, node.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		node.GetLogger().Error(err.Error())
	}
	return nil
}

func (app *App) invertGrid(node *Node.Node, args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = 1 - app.grid[row][col]
		}
	}
	err := node.AsyncMessage(topics.PROPGATE_GRID, node.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		node.GetLogger().Error(err.Error())
	}
	return nil
}

func (app *App) chessGrid(node *Node.Node, args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = (row + col) % 2
		}
	}
	err := node.AsyncMessage(topics.PROPGATE_GRID, node.GetName(), dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		node.GetLogger().Error(err.Error())
	}
	return nil
}
