package appGameOfLife

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Randomizer"
	"sync"
)

type App struct {
	client *Client.Client

	randomizer *Randomizer.Randomizer
	grid       [][]int
	mutex      sync.Mutex
	gridRows   int
	gridCols   int
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	app := &App{
		randomizer: Randomizer.New(Randomizer.GetSystemTime()),
		client:     client,

		grid:     nil,
		gridRows: 90,
		gridCols: 140,
	}
	return app, nil
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
