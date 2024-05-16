package appGameOfLife

import (
	"Systemge/Error"
	"Systemge/Message"
)

func (app *App) MessageHandler(message *Message.Message) error {
	switch message.Type {
	case "gridChange":
		app.mutex.Lock()
		defer app.mutex.Unlock()
		gridChange := UnmarshalGridChange(message.Body)
		app.grid[gridChange.Row][gridChange.Column] = gridChange.State

		app.messageBrokerClient.Send(Message.New("getGridChange", app.name, gridChange.Marshal()))
	case "getGridUnicast":
		app.mutex.Lock()
		defer app.mutex.Unlock()
		app.messageBrokerClient.Send(Message.New("websocketUnicast", app.name, message.Body+"|"+gridToString(app.grid)))
	default:
		return Error.New("Unknown message type: " + message.Type)
	}
	return nil
}
