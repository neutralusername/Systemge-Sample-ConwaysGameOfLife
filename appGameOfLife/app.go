package appGameOfLife

import (
	"sync"

	"github.com/neutralusername/systemge/configs"
	"github.com/neutralusername/systemge/connectionChannel"
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
		gridRows: 90,
		gridCols: 140,
		toroidal: true,

		listener: channelListener,
	}
	grid := make([][]int, app.gridRows)
	for i := range grid {
		grid[i] = make([]int, app.gridCols)
	}
	app.grid = grid

	channelAccepter, err := serviceAccepter.New(
		channelListener,
		&configs.Accepter{},
		&configs.Routine{},
		func(connection systemge.Connection[*tools.Message]) error {

			_, err := serviceReader.NewAsync(
				connection,
				&configs.ReaderAsync{},
				&configs.Routine{},
				func(message *tools.Message, connection systemge.Connection[*tools.Message]) {
					// todo
				},
			)
			if err != nil {
				return err
			}

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

	/*
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
					TcpClientConfigs: []*Config.TcpClient{
						{
							Address: "localhost:60001",
						},
					},
					Reconnect:                   true,
					TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
				},
				func(connection SystemgeConnection.SystemgeConnection) error {
					connection.StartMessageHandlingLoop_Sequentially(messageHandler)
					return nil
				},
				func(connection SystemgeConnection.SystemgeConnection) {
					connection.StopMessageHandlingLoop()
				},
			)
			if err := DashboardClientCustomService.New("appGameOfLife_dashboardClient",
				&Config.DashboardClient{
					TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
					TcpClientConfig: &Config.TcpClient{
						Address: "localhost:60000",
					},
				},
				app.systemgeClient,
				Commands.Handlers{
					"randomize":      app.randomizeGrid,
					"invert":         app.invertGrid,
					"chess":          app.chessGrid,
					"toggleToroidal": app.toggleToroidal,
				}).Start(); err != nil {
				panic(err)
			}
			return app
	*/
	return nil
}

/*
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
*/
