package appGameOfLife

import (
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"

	"github.com/neutralusername/systemge/configs"
	"github.com/neutralusername/systemge/helpers"
	"github.com/neutralusername/systemge/listenerChannel"
	"github.com/neutralusername/systemge/listenerTcp"
	"github.com/neutralusername/systemge/server"
	"github.com/neutralusername/systemge/systemge"
	"github.com/neutralusername/systemge/tools"
	"github.com/neutralusername/systemge/typedListener"
)

type App struct {
	grid     [][]int
	mutex    sync.RWMutex
	gridRows int
	gridCols int
	toroidal bool

	Listener   systemge.Listener[*tools.Message, systemge.Connection[*tools.Message]]
	connection systemge.Connection[*tools.Message]
}

func NewChannel() systemge.Listener[*tools.Message, systemge.Connection[*tools.Message]] {
	listener, err := listenerChannel.New[*tools.Message](
		"listenerTcp",
	)
	if err != nil {
		panic(err)
	}
	if err := listener.Start(); err != nil {
		panic(err)
	}
	return listener
}

func NewTcpListener() systemge.Listener[*tools.Message, systemge.Connection[*tools.Message]] {
	listener, err := listenerTcp.New(
		"listenerTcp",
		&configs.TcpListener{
			Port:   60001,
			Domain: "localhost",
		},
		&configs.TcpBufferedReader{},
	)
	if err != nil {
		panic(err)
	}
	typedListener, err := typedListener.New(
		listener,
		tools.SerializeMessage,
		tools.DeserializeMessage,
	)
	if err != nil {
		panic(err)
	}
	if err := listener.Start(); err != nil {
		panic(err)
	}
	return typedListener
}

func NewApp(listener systemge.Listener[*tools.Message, systemge.Connection[*tools.Message]]) *App {
	app := &App{
		grid:     nil,
		gridRows: 50,
		gridCols: 100,
		toroidal: true,

		Listener: listener,
	}
	grid := make([][]int, app.gridRows)
	for i := range grid {
		grid[i] = make([]int, app.gridCols)
	}
	app.grid = grid

	server, err := server.New(
		listener,
		&configs.Accepter{},
		&configs.Routine{
			MaxConcurrentHandlers: 1,
		},
		&configs.ReaderAsync{},
		&configs.Routine{
			MaxConcurrentHandlers: 1,
		},
		app.acceptHandler,
		app.readHandler,
	)
	if err != nil {
		panic(err)
	}
	if err := server.GetAccepter().GetRoutine().Start(); err != nil {
		panic(err)
	}
	return app
}

func (app *App) acceptHandler(connection systemge.Connection[*tools.Message]) error {
	if app.connection != nil {
		panic("Connection already exists")
	}
	app.connection = connection
	return nil
}

func (app *App) readHandler(message *tools.Message, connection systemge.Connection[*tools.Message]) {
	switch message.GetTopic() {
	case topics.GRID_CHANGE:
		gridChange := dto.UnmarshalGridChange(message.GetPayload())
		app.grid[gridChange.Row][gridChange.Column] = gridChange.State
		app.propagateBoard()

	case topics.NEXT_GENERATION:
		app.calcNextGeneration()
		app.propagateBoard()

	case topics.SET_GRID:
		app.mutex.Lock()
		defer app.mutex.Unlock()

		if len(message.GetPayload()) != app.gridCols*app.gridRows {
			return
		}
		for row := 0; row < app.gridRows; row++ {
			for col := 0; col < app.gridCols; col++ {
				app.grid[row][col] = helpers.StringToInt(string(message.GetPayload()[row*app.gridCols+col]))
			}
		}
		app.propagateBoard()

	case topics.GET_GRID:
		if err := connection.Write(
			tools.NewMessage(
				topics.GET_GRID,
				dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(),
				message.GetSyncToken(),
				true,
			),
			0,
		); err != nil {
			panic(err)
		}
	}
}

func (app *App) propagateBoard() {
	if app.connection == nil {
		return
	}
	if err := app.connection.Write(
		tools.NewMessage(
			topics.PROPAGATE_GRID,
			dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(),
			"",
			false,
		),
		0,
	); err != nil {
		panic(err)
	}
}

func (app *App) calcNextGeneration() {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	nextGrid := make([][]int, app.gridRows)
	for i := range nextGrid {
		nextGrid[i] = make([]int, app.gridCols)
	}
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			aliveNeighbours := 0
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					if app.toroidal {
						neighbourRow := (row + i + app.gridRows) % app.gridRows
						neighbourCol := (col + j + app.gridCols) % app.gridCols
						aliveNeighbours += app.grid[neighbourRow][neighbourCol]
					} else {
						neighbourRow := row + i
						neighbourCol := col + j
						if neighbourRow >= 0 && neighbourRow < app.gridRows && neighbourCol >= 0 && neighbourCol < app.gridCols {
							aliveNeighbours += app.grid[neighbourRow][neighbourCol]
						}
					}
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

func (app *App) toggleToroidal(args []string) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	app.toroidal = !app.toroidal
	return "sucess", nil
}

func (app *App) randomizeGrid(args []string) (string, error) {
	percentageOfAliveCells := int64(50)
	if len(args) > 0 {
		percentageOfAliveCells = helpers.StringToInt64(args[0])
	}
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			if tools.GenerateRandomNumber(1, 100) <= percentageOfAliveCells {
				app.grid[row][col] = 1
			} else {
				app.grid[row][col] = 0
			}
		}
	}
	return "success", nil
}

func (app *App) invertGrid(args []string) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = 1 - app.grid[row][col]
		}
	}
	return "success", nil
}

func (app *App) chessGrid(args []string) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = (row + col) % 2
		}
	}
	return "success", nil
}
