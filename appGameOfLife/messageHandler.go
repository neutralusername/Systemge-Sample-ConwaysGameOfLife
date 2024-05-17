package appGameOfLife

import (
	"Systemge/Message"
)

func (app *App) GridChange(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	gridChange := UnmarshalGridChange(message.Body)
	app.grid[gridChange.Row][gridChange.Column] = gridChange.State
	app.messageBrokerClient.Send(Message.New("getGridChange", app.name, gridChange.Marshal()))
	return nil
}

func (app *App) GetGridUnicast(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.messageBrokerClient.Send(Message.New("websocketUnicast", app.name, message.Body+"|"+"getGrid"+"|"+gridToString(app.grid)))
	return nil
}
