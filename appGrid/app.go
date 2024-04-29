package appGrid

import (
	"Systemge/ApplicationServer"
	"sync"
)

type App struct {
	grid      [10][10]bool
	mutex     sync.Mutex
	appServer *ApplicationServer.Server
}

func Create(appServer *ApplicationServer.Server) ApplicationServer.Application {
	return &App{
		grid:      [10][10]bool{},
		appServer: appServer,
	}
}

func gridToString(grid [10][10]bool) string {
	gridString := ""
	for row := 0; row < 10; row++ {
		for col := 0; col < 10; col++ {
			if grid[row][col] {
				gridString += "1"
			} else {
				gridString += "0"
			}
		}
	}
	return gridString
}
