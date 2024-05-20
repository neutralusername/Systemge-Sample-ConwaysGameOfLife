package appGameOfLife

import (
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"sync"
)

type App struct {
	name                string
	grid                [GRIDROWS][GRIDCOLS]int
	mutex               sync.Mutex
	logger              *Utilities.Logger
	messageBrokerClient *MessageBrokerClient.Client
}

func New(name string, logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client) *App {
	app := &App{
		name:                name,
		grid:                [GRIDROWS][GRIDCOLS]int{},
		logger:              logger,
		messageBrokerClient: messageBrokerClient,
	}
	go app.calcNextGeneration()
	return app
}

func (app *App) calcNextGeneration() {
	nextGrid := [GRIDROWS][GRIDCOLS]int{}
	for row := 0; row < GRIDROWS; row++ {
		for col := 0; col < GRIDCOLS; col++ {
			aliveNeighbours := 0
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					neighbourRow := (row + i + GRIDROWS) % GRIDROWS
					neighbourCol := (col + j + GRIDCOLS) % GRIDCOLS
					aliveNeighbours += app.grid[neighbourRow][neighbourCol]
				}
			}
			aliveNeighbours -= app.grid[row][col]
			if app.grid[row][col] == 1 && (aliveNeighbours < 2 || aliveNeighbours > 3) {
				nextGrid[row][col] = 0
			} else if app.grid[row][col] == 0 && aliveNeighbours == 3 {
				nextGrid[row][col] = 1
			} else {
				nextGrid[row][col] = app.grid[row][col]
			}
		}
	}
	app.grid = nextGrid
}
