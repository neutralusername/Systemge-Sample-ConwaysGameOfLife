package appGameOfLife

import (
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topics"
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Node"
)

type App struct {
	grid                 [][]int
	mutex                sync.Mutex
	gridRows             int
	gridCols             int
	toroidal             bool
	commandHandlers      map[string]Node.CommandHandler
	syncMessageHandlers  map[string]Node.SyncMessageHandler
	asyncMessageHandlers map[string]Node.AsyncMessageHandler
}

var SYSTEMGE_CONFIG = &Config.Systemge{
	HandleMessagesSequentially: false,

	SyncRequestTimeoutMs:            10000,
	TcpTimeoutMs:                    5000,
	MaxConnectionAttempts:           0,
	ConnectionAttemptDelayMs:        1000,
	StopAfterOutgoingConnectionLoss: true,
	ServerConfig: &Config.TcpServer{
		Port: 60001,
	},
	EndpointConfigs: []*Config.TcpEndpoint{
		{
			Address: "localhost:60002",
		},
	},
	IncomingMessageByteLimit: 0,
	MaxPayloadSize:           0,
	MaxTopicSize:             0,
	MaxSyncTokenSize:         0,
}

func New() *App {
	app := &App{
		grid:     nil,
		gridRows: 90,
		gridCols: 140,
		toroidal: true,
	}
	app.commandHandlers = map[string]Node.CommandHandler{
		"randomize":      app.randomizeGrid,
		"invert":         app.invertGrid,
		"chess":          app.chessGrid,
		"toggleToroidal": app.toggleToroidal,
	}
	app.syncMessageHandlers = map[string]Node.SyncMessageHandler{
		topics.GET_GRID: app.getGridSync,
	}
	app.asyncMessageHandlers = map[string]Node.AsyncMessageHandler{
		topics.GRID_CHANGE:     app.gridChange,
		topics.NEXT_GENERATION: app.nextGeneration,
		topics.SET_GRID:        app.setGrid,
	}
	return app
}

func (app *App) GetCommandHandlers() map[string]Node.CommandHandler {
	return app.commandHandlers
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
