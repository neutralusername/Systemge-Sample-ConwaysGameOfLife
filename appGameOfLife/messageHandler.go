package appGameOfLife

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/TypeDefinition"
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

		getGridChangeMsg := typeDefinitions.GET_GRID_CHANGE.New([]string{message.Payload[0][0]}, []string{message.Payload[1][0]}, []string{Utilities.BoolToString(newGridState)})
		app.messageBroker.Send(TypeDefinition.NewWSPropagateMessage([]string{}, getGridChangeMsg))
	case typeDefinitions.REQUEST_GRID.Name:
		app.mutex.Lock()
		defer app.mutex.Unlock()
		getGridMsg := typeDefinitions.GET_GRID.New([]string{gridToString(app.grid)})
		app.messageBroker.Send(TypeDefinition.NewWSPropagateMessage(message.Payload[0], getGridMsg))
	default:
		return Error.New("Unknown message type: " + message.TypeName)
	}
	return nil
}
