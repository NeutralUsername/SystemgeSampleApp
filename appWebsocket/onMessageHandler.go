package appWebsocket

import (
	"Systemge/Message"
	"Systemge/Websocket"
)

func (app *App) OnMessageHandler(connection *Websocket.Connection, message *Message.Message) {
	switch message.Type {
	case "heartbeat":
		connection.ResetWatchdog()
	default:
		message.Origin = connection.Id
		message.SyncKey = ""
		err := app.messageBrokerClient.AsyncMessage(message)
		if err != nil {
			connection.Send([]byte(Message.New("error", "", app.name, err.Error()).Serialize()))
		}
	}
}
