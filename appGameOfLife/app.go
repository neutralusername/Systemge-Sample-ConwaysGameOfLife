package appGameOfLife

import (
	"Systemge/Node"
	"Systemge/Utilities"
	"sync"
)

type App struct {
	randomizer *Utilities.Randomizer
	grid       [][]int
	mutex      sync.Mutex
	gridRows   int
	gridCols   int
	toroidal   bool
}

func New() Node.Application {
	app := &App{
		randomizer: Utilities.NewRandomizer(Utilities.GetSystemTime()),

		grid:     nil,
		gridRows: 90,
		gridCols: 140,
		toroidal: true,
	}
	return app
}

func (app *App) OnStart(client *Node.Node) error {
	grid := make([][]int, app.gridRows)
	for i := range grid {
		grid[i] = make([]int, app.gridCols)
	}
	app.grid = grid
	return nil
}

func (app *App) OnStop(client *Node.Node) error {
	return nil
}
