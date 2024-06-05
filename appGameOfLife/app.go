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
	messageBrokerClient *MessageBrokerClient.Client

	randomizer *Utilities.Randomizer
	grid       [][]int
	mutex      sync.Mutex
	gridRows   int
	gridCols   int
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
