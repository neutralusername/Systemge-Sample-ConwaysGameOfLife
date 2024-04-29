package appGrid

import (
	"Systemge/ApplicationServer"
	"Systemge/RequestServer"
	"SystemgeSampleApp/typeDefinitions"
	"sync"
	"time"
)

const GRIDSIZE = 75

type App struct {
	grid      [GRIDSIZE][GRIDSIZE]bool
	mutex     sync.Mutex
	appServer *ApplicationServer.Server

	websocketEndpoint RequestServer.Endpoint
}

func Create(appServer *ApplicationServer.Server, websocketEndpoint RequestServer.Endpoint) ApplicationServer.Application {
	app := &App{
		grid:      [GRIDSIZE][GRIDSIZE]bool{},
		appServer: appServer,

		websocketEndpoint: websocketEndpoint,
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
	app.websocketEndpoint.Request(typeDefinitions.PROPAGATE_GRID_REQUEST.Create([]string{gridToString(app.grid)}))
	time.Sleep(5 * time.Second)
	app.calcNextGeneration()
}
