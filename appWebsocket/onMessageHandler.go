package appWebsocket

import (
	"Systemge/Message"
	"Systemge/Websocket"
)

func (app *App) OnMessageHandler(connection *Websocket.Connection, message *Message.Message) {
	message.Origin = connection.Id
	message.SyncKey = ""
	switch message.Type {
	case "heartbeat":
		connection.ResetWatchdog()
	case "gridChange":
		err := app.messageBrokerClient.AsyncMessage(message)
		if err != nil {
			connection.Send([]byte(Message.New("error", "", app.name, err.Error()).Serialize()))
		}
	case "nextGeneration":
		err := app.messageBrokerClient.AsyncMessage(message)
		if err != nil {
			connection.Send([]byte(Message.New("error", "", app.name, err.Error()).Serialize()))
		}
	default:
		connection.Send([]byte(Message.New("error", "", app.name, "Unknown message type").Serialize()))
	}
}