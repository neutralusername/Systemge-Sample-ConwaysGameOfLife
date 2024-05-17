package appGameOfLife

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"sync"
	"time"
)

type App struct {
	name                string
	grid                [GRIDSIZE][GRIDSIZE]bool
	mutex               sync.Mutex
	logger              *Utilities.Logger
	messageBrokerClient *MessageBrokerClient.Client
}

func New(name string, logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client) *App {
	app := &App{
		name:                name,
		grid:                [GRIDSIZE][GRIDSIZE]bool{},
		logger:              logger,
		messageBrokerClient: messageBrokerClient,
	}
	go app.calcNextGeneration()
	return app
}

func gridToString(grid [GRIDSIZE][GRIDSIZE]bool) string {
	gridString := ""
	for row := 0; row < GRIDSIZE; row++ {
		for col := 0; col < GRIDSIZE; col++ {
			if grid[row][col] {
				gridString += "1"
			} else {
				gridString += "0"
			}
		}
	}
	return gridString
}

func (app *App) calcNextGeneration() {
	nextGrid := [GRIDSIZE][GRIDSIZE]bool{}
	for row := 0; row < GRIDSIZE; row++ {
		for col := 0; col < GRIDSIZE; col++ {
			aliveNeighbors := 0
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					if app.grid[(row+i+GRIDSIZE)%GRIDSIZE][(col+j+GRIDSIZE)%GRIDSIZE] {
						aliveNeighbors++
					}
				}
			}
			if app.grid[row][col] {
				aliveNeighbors--
			}
			if app.grid[row][col] && aliveNeighbors < 2 {
				nextGrid[row][col] = false
			}
			if app.grid[row][col] && (aliveNeighbors == 2 || aliveNeighbors == 3) {
				nextGrid[row][col] = true
			}
			if app.grid[row][col] && aliveNeighbors > 3 {
				nextGrid[row][col] = false
			}
			if !app.grid[row][col] && aliveNeighbors == 3 {
				nextGrid[row][col] = true
			}
		}
	}

	app.grid = nextGrid
	err := app.messageBrokerClient.Send(Message.New("getGrid", app.name, gridToString(app.grid)))
	if err != nil {
		app.logger.Log(Error.New(err.Error()).Error())
	}
	time.Sleep(5 * time.Second)
	app.calcNextGeneration()
}
