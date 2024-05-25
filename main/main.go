package main

import (
	"Systemge/HTTP"
	"Systemge/MessageBrokerClient"
	"Systemge/MessageBrokerServer"
	"Systemge/TopicResolutionServer"
	"Systemge/Utilities"
	"Systemge/Websocket"
	"SystemgeSampleApp/appGameOfLife"
	"SystemgeSampleApp/appWebsocket"
	"time"
)

const MESSAGEBROKERSERVER_A_ADDRESS = ":60003"
const MESSAGEBROKERSERVER_B_ADDRESS = ":60004"
const TOPICRESOLUTIONSERVER_ADDRESS = ":60002"
const HTTP_DEV_PORT = ":8080"
const WEBSOCKET_PORT = ":8443"

func main() {
	logger := Utilities.NewLogger("error_log.txt")

	messageBrokerServerA := MessageBrokerServer.New("messageBrokerServer", MESSAGEBROKERSERVER_A_ADDRESS, logger)
	messageBrokerServerA.Start()
	messageBrokerServerA.AddTopic("getGridSync")
	messageBrokerServerA.AddTopic("gridChange")
	messageBrokerServerA.AddTopic("nextGeneration")
	messageBrokerServerA.AddTopic("setGrid")

	messageBrokerServerB := MessageBrokerServer.New("messageBrokerServer", MESSAGEBROKERSERVER_B_ADDRESS, logger)
	messageBrokerServerB.Start()
	messageBrokerServerB.AddTopic("getGrid")
	messageBrokerServerB.AddTopic("getGridChange")

	topicResolutionServer := TopicResolutionServer.New("topicResolutionServer", TOPICRESOLUTIONSERVER_ADDRESS, logger)
	topicResolutionServer.Start()
	topicResolutionServer.RegisterTopic("getGridSync", MESSAGEBROKERSERVER_A_ADDRESS)
	topicResolutionServer.RegisterTopic("gridChange", MESSAGEBROKERSERVER_A_ADDRESS)
	topicResolutionServer.RegisterTopic("nextGeneration", MESSAGEBROKERSERVER_A_ADDRESS)
	topicResolutionServer.RegisterTopic("setGrid", MESSAGEBROKERSERVER_A_ADDRESS)
	topicResolutionServer.RegisterTopic("getGrid", MESSAGEBROKERSERVER_B_ADDRESS)
	topicResolutionServer.RegisterTopic("getGridChange", MESSAGEBROKERSERVER_B_ADDRESS)

	messageBrokerClientWebsocket := MessageBrokerClient.New("messageBrokerClientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, logger)
	messageBrokerClientWebsocket.Connect()

	messageBrokerClientGameOfLife := MessageBrokerClient.New("messageBrokerClientGrid", TOPICRESOLUTIONSERVER_ADDRESS, logger)
	messageBrokerClientGameOfLife.Connect()

	websocketServer := Websocket.New("websocketServer", messageBrokerClientWebsocket)
	appWebsocket := appWebsocket.New("websocketApp", logger, messageBrokerClientWebsocket, websocketServer)
	appGameOfLife := appGameOfLife.New("gameOfLifeApp", logger, messageBrokerClientGameOfLife, 90, 140)

	messageBrokerClientGameOfLife.RegisterIncomingSyncTopic("getGridSync", appGameOfLife.GetGridSync)
	messageBrokerClientGameOfLife.RegisterIncomingAsyncTopic("gridChange", appGameOfLife.GridChange)
	messageBrokerClientGameOfLife.RegisterIncomingAsyncTopic("nextGeneration", appGameOfLife.NextGeneration)
	messageBrokerClientGameOfLife.RegisterIncomingAsyncTopic("setGrid", appGameOfLife.SetGrid)

	messageBrokerClientWebsocket.RegisterIncomingAsyncTopic("getGrid", appWebsocket.GetGrid)
	messageBrokerClientWebsocket.RegisterIncomingAsyncTopic("getGridChange", appWebsocket.GetGridChange)

	websocketServer.SetOnMessageHandler(appWebsocket.OnMessageHandler)
	websocketServer.SetOnConnectHandler(appWebsocket.OnConnectHandler)

	websocketServer.Start()

	HTTPServerServe := HTTP.New(HTTP_DEV_PORT, "HTTPfrontend", false, "", "")
	HTTPServerServe.RegisterPattern("/", HTTP.SendDirectory("../frontend"))
	HTTPServerServe.Start()

	HTTPServerWebsocket := HTTP.New(WEBSOCKET_PORT, "HTTPwebsocket", false, "", "")
	HTTPServerWebsocket.RegisterPattern("/ws", HTTP.PromoteToWebsocket(websocketServer))
	HTTPServerWebsocket.Start()

	time.Sleep(1000000 * time.Second)
}
