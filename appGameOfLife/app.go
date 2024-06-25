package appGameOfLife

import (
	"Systemge/Client"
	"Systemge/Randomizer"
	"sync"
)

type App struct {
	randomizer *Randomizer.Randomizer
	grid       [][]int
	mutex      sync.Mutex
	gridRows   int
	gridCols   int
	toroidal   bool
}

func New() Client.Application {
	app := &App{
		randomizer: Randomizer.New(Randomizer.GetSystemTime()),

		grid:     nil,
		gridRows: 90,
		gridCols: 140,
		toroidal: true,
	}
	return app
}

func (app *App) OnStart(client *Client.Client) error {
	grid := make([][]int, app.gridRows)
	for i := range grid {
		grid[i] = make([]int, app.gridCols)
	}
	app.grid = grid
	return nil
}

func (app *App) OnStop(client *Client.Client) error {
	return nil
}
