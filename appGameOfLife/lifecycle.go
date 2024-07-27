package appGameOfLife

import "github.com/neutralusername/Systemge/Node"

func (app *App) OnStart(node *Node.Node) error {
	grid := make([][]int, app.gridRows)
	for i := range grid {
		grid[i] = make([]int, app.gridCols)
	}
	app.grid = grid
	return nil
}
