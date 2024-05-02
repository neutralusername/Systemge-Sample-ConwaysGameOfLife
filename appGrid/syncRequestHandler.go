package appGrid

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/TypeDefinition"
	"Systemge/Utilities"
	"SystemgeSampleApp/typeDefinitions"
)

func (app *App) SyncRequestHandler(message *Message.Message) *Message.Message {
	switch message.TypeName {
	case typeDefinitions.SET_GRID_REQUEST.Name:
		app.mutex.Lock()
		defer app.mutex.Unlock()
		newGridState := !app.grid[Utilities.StringToInt(message.Payload[0][0])][Utilities.StringToInt(message.Payload[1][0])]
		app.grid[Utilities.StringToInt(message.Payload[0][0])][Utilities.StringToInt(message.Payload[1][0])] = newGridState
		return typeDefinitions.SET_GRID_REQUEST.Response.New([]string{Utilities.BoolToString(newGridState)})
	case typeDefinitions.GET_GRID_REQUEST.Name:
		app.mutex.Lock()
		defer app.mutex.Unlock()
		return typeDefinitions.GET_GRID_REQUEST.Response.New([]string{gridToString(app.grid)})
	default:
		return TypeDefinition.NewErrorMessage(Error.New("Unknown message type: " + message.TypeName))
	}
}
