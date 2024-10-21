package appGameOfLife

import (
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"

	"github.com/neutralusername/systemge/configs"
	"github.com/neutralusername/systemge/connectionChannel"
	"github.com/neutralusername/systemge/helpers"
	"github.com/neutralusername/systemge/listenerChannel"
	"github.com/neutralusername/systemge/serviceAccepter"
	"github.com/neutralusername/systemge/serviceReader"
	"github.com/neutralusername/systemge/systemge"
	"github.com/neutralusername/systemge/tools"
)

type App struct {
	grid     [][]int
	mutex    sync.RWMutex
	gridRows int
	gridCols int
	toroidal bool

	listener    systemge.Listener[*tools.Message, systemge.Connection[*tools.Message]]
	accepter    *serviceAccepter.Accepter[*tools.Message]
	connections map[systemge.Connection[*tools.Message]]struct{}
}

var ConnectionChannel chan<- *connectionChannel.ConnectionRequest[*tools.Message]

func New() *App {
	channelListener, err := listenerChannel.New[*tools.Message]("listenerChannel")
	if err != nil {
		panic(err)
	}
	ConnectionChannel = channelListener.(*listenerChannel.ChannelListener[*tools.Message]).GetConnectionChannel() // this should be less complicated (make a function that takes systemgeListener and returns either err or this channel)

	app := &App{
		grid:     nil,
		gridRows: 50,
		gridCols: 100,
		toroidal: true,

		connections: make(map[systemge.Connection[*tools.Message]]struct{}),
		listener:    channelListener,
	}
	grid := make([][]int, app.gridRows)
	for i := range grid {
		grid[i] = make([]int, app.gridCols)
	}
	app.grid = grid

	channelAccepter, err := serviceAccepter.New(
		channelListener,
		&configs.Accepter{},
		&configs.Routine{
			MaxConcurrentHandlers: 1,
		},
		func(connection systemge.Connection[*tools.Message]) error {

			reader, err := serviceReader.NewAsync(
				connection,
				&configs.ReaderAsync{},
				&configs.Routine{
					MaxConcurrentHandlers: 10,
				},
				func(message *tools.Message, connection systemge.Connection[*tools.Message]) {

					switch message.GetTopic() {
					case topics.GRID_CHANGE:
						gridChange := dto.UnmarshalGridChange(message.GetPayload())
						app.grid[gridChange.Row][gridChange.Column] = gridChange.State
						err := connection.Write(
							tools.NewMessage(
								topics.PROPAGATE_GRID_CHANGE,
								message.GetPayload(),
								"",
								false,
							),
							0,
						)
						if err != nil {
							panic(err)
						}

					case topics.NEXT_GENERATION:
						app.calcNextGeneration()
						err := connection.Write(
							tools.NewMessage(
								topics.PROPAGATE_GRID,
								dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(),
								"",
								false,
							),
							0,
						)
						if err != nil {
							panic(err)
						}

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
						err := connection.Write(
							tools.NewMessage(
								topics.PROPAGATE_GRID,
								dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(),
								"",
								false,
							),
							0,
						)
						if err != nil {
							panic(err)
						}

					case topics.GET_GRID:
						err := connection.Write(
							tools.NewMessage(
								topics.GET_GRID,
								dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(),
								message.GetSyncToken(),
								true,
							),
							0,
						)
						if err != nil {
							panic(err)
						}
					}
				},
			)
			if err != nil {
				return err
			}
			reader.GetRoutine().Start()

			app.mutex.Lock()
			app.connections[connection] = struct{}{}
			app.mutex.Unlock()

			go func() { // abstract on close handler
				<-connection.GetCloseChannel()

				app.mutex.Lock()
				defer app.mutex.Unlock()

				delete(app.connections, connection)
			}()

			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	app.accepter = channelAccepter
	app.accepter.GetRoutine().Start()

	return app
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

/*
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
*/
