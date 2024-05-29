package appGameOfLife

import (
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"sync"
)

type App struct {
	messageBrokerClient *MessageBrokerClient.Client
	name                string
	grid                [][]int
	mutex               sync.Mutex
	gridRows            int
	gridCols            int
	logger              *Utilities.Logger
	randomizer          *Utilities.Randomizer
}

func New(name string, logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client, gridRows, gridCols int) *App {
	grid := make([][]int, gridRows)
	for i := range grid {
		grid[i] = make([]int, gridCols)
	}
	app := &App{
		messageBrokerClient: messageBrokerClient,
		name:                name,
		grid:                grid,
		gridRows:            gridRows,
		gridCols:            gridCols,
		logger:              logger,
		randomizer:          Utilities.NewRandomizer(Utilities.GetSystemTime()),
	}
	return app
}

func (app *App) calcNextGeneration() {
	nextGrid := make([][]int, app.gridRows)
	for i := range nextGrid {
		nextGrid[i] = make([]int, app.gridCols)
	}
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			aliveNeighbours := 0
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					neighbourRow := (row + i + app.gridRows) % app.gridRows
					neighbourCol := (col + j + app.gridCols) % app.gridCols
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

func (app *App) RandomizeGrid() {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = app.randomizer.GenerateRandomNumber(0, 1)
		}
	}
	err := app.messageBrokerClient.AsyncMessage(Message.NewAsync("getGrid", app.name, newGrid(app.grid, app.gridRows, app.gridCols).marshal()))
	if err != nil {
		app.logger.Log(err.Error())
	}
}

func (app *App) InvertGrid() {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = 1 - app.grid[row][col]
		}
	}
	err := app.messageBrokerClient.AsyncMessage(Message.NewAsync("getGrid", app.name, newGrid(app.grid, app.gridRows, app.gridCols).marshal()))
	if err != nil {
		app.logger.Log(err.Error())
	}
}
