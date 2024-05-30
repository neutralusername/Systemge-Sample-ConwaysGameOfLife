package appGameOfLife

import (
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"sync"
)

const gridRows = 90
const gridCols = 140

type App struct {
	logger              *Utilities.Logger
	randomizer          *Utilities.Randomizer
	messageBrokerClient *MessageBrokerClient.Client

	grid     [][]int
	mutex    sync.Mutex
	gridRows int
	gridCols int
}

func New(logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client) MessageBrokerClient.Application {
	grid := make([][]int, gridRows)
	for i := range grid {
		grid[i] = make([]int, gridCols)
	}
	app := &App{
		logger:              logger,
		randomizer:          Utilities.NewRandomizer(Utilities.GetSystemTime()),
		messageBrokerClient: messageBrokerClient,

		grid:     grid,
		gridRows: gridRows,
		gridCols: gridCols,
	}
	return app
}

func (app *App) GetAsyncMessageHandlers() map[string]MessageBrokerClient.AsyncMessageHandler {
	return map[string]MessageBrokerClient.AsyncMessageHandler{
		"gridChange":     app.GridChange,
		"nextGeneration": app.NextGeneration,
		"setGrid":        app.SetGrid,
	}
}

func (app *App) GetSyncMessageHandlers() map[string]MessageBrokerClient.SyncMessageHandler {
	return map[string]MessageBrokerClient.SyncMessageHandler{
		"getGridSync": app.GetGridSync,
	}
}

func (app *App) GetCustomCommandHandlers() map[string]func() error {
	return map[string]func() error{
		"randomize": app.RandomizeGrid,
		"invert":    app.InvertGrid,
	}
}
