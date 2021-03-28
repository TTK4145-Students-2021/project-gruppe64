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

	cartButtonCh := make(chan hardwareIO.ButtonEvent)
	hallOrderCh := make(chan hardwareIO.ButtonEvent)
	floorArrivalCh := make(chan int)
	timerDurationCh := make(chan float64)
	timedOutCh := make(chan bool)

	go hardwareIO.RunHardware(cartButtonCh, hallOrderCh, floorArrivalCh)
	go timer.RunBlockingTimer(timerDurationCh, timedOutCh)
	go fsm.ElevatorFSM(cartButtonCh, floorArrivalCh, timerDurationCh, timedOutCh)
	for {}

}