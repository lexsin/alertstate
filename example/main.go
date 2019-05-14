package main

import (
	//"fmt"
	"ifmes/modules/alertstate"
)

func init() {
	alertstate.Init(5)
}

func main() {
	alert_state_start()
}

func alert_state_start() {
	alertstate.GLocalCach.Handler = alert_state_handler
	go alertstate.GLocalCach.Start()
	alertstate.StartTest()
}
