package main

import (
	"fmt"
	"ifmes/modules/alertstate"
	"time"
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

func TimeFormUi(stamp int64) string {
	t := time.Unix(stamp, 0)
	tstr := fmt.Sprintf("%d-%d-%d %d:%d:%d",
		t.Year(), int(t.Month()), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return tstr
}
