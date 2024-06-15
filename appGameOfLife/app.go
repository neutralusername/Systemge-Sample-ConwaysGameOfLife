package appGameOfLife

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Utilities"
	"sync"
)

type App struct {
	client *Client.Client

	randomizer *Utilities.Randomizer
	grid       [][]int
	mutex      sync.Mutex
	gridRows   int
	gridCols   int
}

func New(client *Client.Client, args []string) Application.Application {
	app := &App{
		randomizer: Utilities.NewRandomizer(Utilities.GetSystemTime()),
		client:     client,

		grid:     nil,
		gridRows: 90,
		gridCols: 140,
	}
	return app
}

func (app *App) OnStart() error {
	grid := make([][]int, app.gridRows)
	for i := range grid {
		grid[i] = make([]int, app.gridCols)
	}
	app.grid = grid
	return nil
}

func (app *App) OnStop() error {
	return nil
}
