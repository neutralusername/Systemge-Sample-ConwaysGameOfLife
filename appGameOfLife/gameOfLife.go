package appGameOfLife

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
