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
		app.grid[Utilities.StringToInt(message.Payload[0][0])][Utilities.StringToInt(message.Payload[1][0])] = !app.grid[Utilities.StringToInt(message.Payload[0][0])][Utilities.StringToInt(message.Payload[1][0])]
		return typeDefinitions.SET_GRID_REQUEST.Response.Create([]string{gridToString(app.grid)})
	case typeDefinitions.GET_GRID_REQUEST.Name:
		app.mutex.Lock()
		defer app.mutex.Unlock()
		return typeDefinitions.GET_GRID_REQUEST.Response.Create([]string{gridToString(app.grid)})
	default:
		return TypeDefinition.CreateError(Error.Create("Unknown message type: " + message.TypeName))
	}
}
