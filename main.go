package main

import (
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/timer"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	hardwareIO.Init("localhost:15657", hardwareIO.NumFloors)

	buttonEventCh := make(chan hardwareIO.ButtonEvent)
	floorArrivalCh := make(chan int)
	timerDurationCh := make(chan float64)
	timedOutCh := make(chan bool)

	go hardwareIO.RunHardware(buttonEventCh, floorArrivalCh)
	go timer.RunBlockingTimer(timerDurationCh, timedOutCh)
	go fsm.ElevatorFSM(buttonEventCh, floorArrivalCh, timerDurationCh, timedOutCh)
	for {}

}