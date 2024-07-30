package appGameOfLife

import (
	"SystemgeSampleConwaysGameOfLife/dto"
	"SystemgeSampleConwaysGameOfLife/topics"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) GetSystemgeComponentConfig() *Config.Systemge {
	return &Config.Systemge{
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
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return app.syncMessageHandlers
}

func (app *App) getGridSync(node *Node.Node, message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	return dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal(), nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return app.asyncMessageHandlers
}

func (app *App) gridChange(node *Node.Node, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	gridChange := dto.UnmarshalGridChange(message.GetPayload())
	app.grid[gridChange.Row][gridChange.Column] = gridChange.State
	node.AsyncMessage(topics.PROPAGATE_GRID_CHANGE, gridChange.Marshal())
	return nil
}

func (app *App) nextGeneration(node *Node.Node, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.calcNextGeneration()
	err := node.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to propagate grid: " + err.Error())
		}
	}
	return nil
}

func (app *App) setGrid(node *Node.Node, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if len(message.GetPayload()) != app.gridCols*app.gridRows {
		return Error.New("Invalid grid size", nil)
	}
	for row := 0; row < app.gridRows; row++ {
		for col := 0; col < app.gridCols; col++ {
			app.grid[row][col] = Helpers.StringToInt(string(message.GetPayload()[row*app.gridCols+col]))
		}
	}
	err := node.AsyncMessage(topics.PROPGATE_GRID, dto.NewGrid(app.grid, app.gridRows, app.gridCols).Marshal())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to propagate grid: " + err.Error())
		}
	}
	return nil
}
