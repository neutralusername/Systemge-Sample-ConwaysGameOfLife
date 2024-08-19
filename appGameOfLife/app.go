package appGameOfLife

import (
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeClient"
	"github.com/neutralusername/Systemge/SystemgeMessageHandler"
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
	app.systemgeClient = SystemgeClient.New(&Config.SystemgeClient{
		Name: "systemgeClent",
		EndpointConfigs: []*Config.TcpEndpoint{
			{
				Address: "localhost:60001",
			},
		},
		ConnectionConfig: &Config.SystemgeConnection{},
	}, &Config.SystemgeReceiver{}, SystemgeMessageHandler.New(SystemgeMessageHandler.AsyncMessageHandlers{
		topics.GRID_CHANGE:     app.gridChange,
		topics.NEXT_GENERATION: app.nextGeneration,
		topics.SET_GRID:        app.setGrid,
	}, SystemgeMessageHandler.SyncMessageHandlers{
		topics.GET_GRID: app.getGridSync,
	}))
	if err := app.systemgeClient.Start(); err != nil {
		panic(err)
	}
	return app
}

func (app *App) getGridSync(message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	return dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(), nil
}

func (app *App) gridChange(message *Message.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	gridChange := dto.UnmarshalGridChange(message.GetPayload())
	app.grid[gridChange.Row][gridChange.Column] = gridChange.State
	app.systemgeClient.AsyncMessage(topics.PROPAGATE_GRID_CHANGE, gridChange.Marshal())
}

func (app *App) nextGeneration(message *Message.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.calcNextGeneration()
	app.systemgeClient.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
}

func (app *App) setGrid(message *Message.Message) {
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

/*
	app.commandHandlers = map[string]Node.CommandHandler{
		"randomize":      app.randomizeGrid,
		"invert":         app.invertGrid,
		"chess":          app.chessGrid,
		"toggleToroidal": app.toggleToroidal,
	}


	func (app *App) toggleToroidal(node *Node.Node, args []string) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.toroidal = !app.toroidal
	return "sucess", nil
}

func (app *App) randomizeGrid(node *Node.Node, args []string) (string, error) {
	percentageOfAliveCells := int64(50)
	if len(args) > 0 {
		percentageOfAliveCells = Helpers.StringToInt64(args[0])
	}
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			if node.GetRandomizer().GenerateRandomNumber(1, 100) <= percentageOfAliveCells {
				app.grid[row][col] = 1
			} else {
				app.grid[row][col] = 0
			}
		}
	}
	err := node.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to propagate grid: " + err.Error())
		}
	}
	return "success", nil
}

func (app *App) invertGrid(node *Node.Node, args []string) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = 1 - app.grid[row][col]
		}
	}
	err := node.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to propagate grid: " + err.Error())
		}
	}
	return "success", nil
}

func (app *App) chessGrid(node *Node.Node, args []string) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = (row + col) % 2
		}
	}
	err := node.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to propagate grid: " + err.Error())
		}
	}
	return "success", nil
}

*/
