package appGameOfLife

import (
	"Systemge/Error"
	"Systemge/Message"
)

func (app *App) GridChange(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	gridChange := UnmarshalGridChange(message.Body)
	app.grid[gridChange.Row][gridChange.Column] = gridChange.State
	app.messageBrokerClient.AsyncMessage(Message.New("getGridChange", app.name, "", gridChange.Marshal()))
	return nil
}

func (app *App) GetGridSync(message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	return gridToString(app.grid), nil
}

func (app *App) NextGeneration(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.calcNextGeneration()
	err := app.messageBrokerClient.AsyncMessage(Message.New("getGrid", app.name, "", gridToString(app.grid)))
	if err != nil {
		app.logger.Log(Error.New(err.Error()).Error())
	}
	return nil
}
