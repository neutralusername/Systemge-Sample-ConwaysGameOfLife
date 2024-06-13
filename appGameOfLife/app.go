package appGameOfLife

import (
	"Systemge/Application"
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"sync"
)

type App struct {
	logger              *Utilities.Logger
	messageBrokerClient *MessageBrokerClient.Client

	randomizer *Utilities.Randomizer
	grid       [][]int
	mutex      sync.Mutex
	gridRows   int
	gridCols   int
}

func New(logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client) Application.Application {
	app := &App{
		logger:              logger,
		randomizer:          Utilities.NewRandomizer(Utilities.GetSystemTime()),
		messageBrokerClient: messageBrokerClient,

		grid:     nil,
		gridRows: 90,
		gridCols: 140,
	}
	return app
}

func (app *App) OnStart() error {
	grid := make([][]int, app.gridCols)
	for i := range grid {
		grid[i] = make([]int, app.gridRows)
	}
	app.grid = grid
	return nil
}

func (app *App) OnStop() error {
	return nil
}
