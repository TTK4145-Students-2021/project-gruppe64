package main

import (
	"fmt"
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

	drvButtons := make(chan hardwareIO.ButtonEvent)
	drvFloors  := make(chan int)
	drvObstr   := make(chan bool)
	drvStop    := make(chan bool)

	go hardwareIO.PollButtons(drvButtons)
	go hardwareIO.PollFloorSensor(drvFloors)
	go hardwareIO.PollObstructionSwitch(drvObstr)
	go hardwareIO.PollStopButton(drvStop)

	go timer.RunTimer(timerDurationCh, timedOutCh)
	go fsm.ElevatorFSM(buttonEventCh, floorArrivalCh, timerDurationCh, timedOutCh)

	for {
		select {
		case a := <- drvButtons:
			fmt.Printf("%+v\n", a)
			hardwareIO.SetButtonLamp(a.Button, a.Floor, true)
			buttonEventCh <- a
		case a := <- drvFloors:
			fmt.Printf("%+v\n", a)
			floorArrivalCh <- a
		case a := <- drvObstr:
			fmt.Printf("%+v\n", a)
			if a {
				hardwareIO.SetMotorDirection(hardwareIO.MD_Stop)
			} else {}
		case a := <- drvStop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < hardwareIO.NumFloors; f++ {
				for b := hardwareIO.ButtonType(0); b < 3; b++ {
					hardwareIO.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}