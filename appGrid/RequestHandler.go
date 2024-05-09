package appGrid

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleApp/typeDefinitions"
)

func (app *App) MessageHandler(message *Message.Message) error {
	switch message.TypeName {
	case typeDefinitions.REQUEST_GRID_CHANGE.Name:
		app.mutex.Lock()
		defer app.mutex.Unlock()
		newGridState := !app.grid[Utilities.StringToInt(message.Payload[0][0])][Utilities.StringToInt(message.Payload[1][0])]
		app.grid[Utilities.StringToInt(message.Payload[0][0])][Utilities.StringToInt(message.Payload[1][0])] = newGridState
		app.messageBroker.Message(typeDefinitions.BROADCAST_GRID_CHANGE.New([]string{message.Payload[0][0]}, []string{message.Payload[1][0]}, []string{Utilities.BoolToString(newGridState)}))
	case typeDefinitions.REQUEST_GRID_BROADCAST.Name:
		app.mutex.Lock()
		defer app.mutex.Unlock()
		app.messageBroker.Message(typeDefinitions.BROADCAST_GRID.New([]string{gridToString(app.grid)}))
	case typeDefinitions.REQUEST_GRID_UNICAST.Name:
		app.mutex.Lock()
		defer app.mutex.Unlock()
		app.messageBroker.Message(typeDefinitions.UNICAST_GRID.New([]string{message.Payload[0][0]}, []string{gridToString(app.grid)}))
	default:
		return Error.New("Unknown message type: " + message.TypeName)
	}
	return nil
}
