package appGameOfLife

import (
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"

	"github.com/neutralusername/Systemge/Commands"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeClient"
	"github.com/neutralusername/Systemge/SystemgeConnection"
	"github.com/neutralusername/Systemge/Tools"
)

type App struct {
	grid     [][]int
	mutex    sync.Mutex
	gridRows int
	gridCols int
	toroidal bool

	systemgeClient *SystemgeClient.SystemgeClient
}

func New() *App {
	app := &App{
		grid:     nil,
		gridRows: 90,
		gridCols: 140,
		toroidal: true,
	}
	grid := make([][]int, app.gridRows)
	for i := range grid {
		grid[i] = make([]int, app.gridCols)
	}
	app.grid = grid

	messageHandler := SystemgeConnection.NewTopicExclusiveMessageHandler(
		SystemgeConnection.AsyncMessageHandlers{
			topics.GRID_CHANGE:     app.gridChange,
			topics.NEXT_GENERATION: app.nextGeneration,
			topics.SET_GRID:        app.setGrid,
		},
		SystemgeConnection.SyncMessageHandlers{
			topics.GET_GRID: app.getGridSync,
		},
		nil, nil, 100,
	)
	app.systemgeClient = SystemgeClient.New("appGameOfLife_systemgeClient",
		&Config.SystemgeClient{
			ClientConfigs: []*Config.TcpClient{
				{
					Address: "localhost:60001",
				},
			},
			ConnectionConfig: &Config.TcpSystemgeConnection{},
		},
		func(connection SystemgeConnection.SystemgeConnection) error {
			connection.StartProcessingLoopSequentially(messageHandler)
			return nil
		},
		func(connection SystemgeConnection.SystemgeConnection) {
			connection.StopProcessingLoop()
		},
	)
	if err := Dashboard.NewClient("appGameOfLife_dashboardClient",
		&Config.DashboardClient{
			ConnectionConfig: &Config.TcpSystemgeConnection{},
			ClientConfig: &Config.TcpClient{
				Address: "localhost:60000",
			},
		},
		app.systemgeClient.Start, app.systemgeClient.Stop, app.systemgeClient.GetMetrics, app.systemgeClient.GetStatus,
		Commands.Handlers{
			"randomize":      app.randomizeGrid,
			"invert":         app.invertGrid,
			"chess":          app.chessGrid,
			"toggleToroidal": app.toggleToroidal,
		}).Start(); err != nil {
		panic(err)
	}
	return app
}

func (app *App) getGridSync(connection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	return dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(), nil
}

func (app *App) gridChange(connection SystemgeConnection.SystemgeConnection, message *Message.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	gridChange := dto.UnmarshalGridChange(message.GetPayload())
	app.grid[gridChange.Row][gridChange.Column] = gridChange.State
	app.systemgeClient.AsyncMessage(topics.PROPAGATE_GRID_CHANGE, gridChange.Marshal())
}

func (app *App) nextGeneration(connection SystemgeConnection.SystemgeConnection, message *Message.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.calcNextGeneration()
	app.systemgeClient.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
}

func (app *App) setGrid(connection SystemgeConnection.SystemgeConnection, message *Message.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if len(message.GetPayload()) != app.gridCols*app.gridRows {
		return
	}
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = Helpers.StringToInt(string(message.GetPayload()[row*app.gridCols+col]))
		}
	}
	app.systemgeClient.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
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
		percentageOfAliveCells = Helpers.StringToInt64(args[0])
	}
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			if Tools.GenerateRandomNumber(1, 100) <= percentageOfAliveCells {
				app.grid[row][col] = 1
			} else {
				app.grid[row][col] = 0
			}
		}
	}
	app.systemgeClient.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
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
	app.systemgeClient.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
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
	app.systemgeClient.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	return "success", nil
}
