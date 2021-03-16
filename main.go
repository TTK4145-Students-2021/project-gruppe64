package main

import (
	"fmt"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/timer"
	"runtime"
)

func main()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	numFloors := hardwareIO.NumFloors
	hardwareIO.Init("localhost:15657", numFloors)

	drvButtons := make(chan hardwareIO.ButtonEvent)
	drvFloors  := make(chan int)
	drvObstr   := make(chan bool)
	drvStop    := make(chan bool)

	buttonCh := make(chan hardwareIO.ButtonEvent)
	floorCh := make(chan int)

	elevatorToNOCh := make(chan fsm.Elevator)
	elevatorToFACh := make(chan fsm.Elevator)
	elevatorToOECh := make(chan fsm.Elevator)
	elevatorUpdateCh := make(chan fsm.Elevator)

	runTimerCh := make(chan float64)
	timedOutCh := make(chan bool)

	go hardwareIO.PollButtons(drvButtons)
	go hardwareIO.PollFloorSensor(drvFloors)
	go hardwareIO.PollObstructionSwitch(drvObstr)
	go hardwareIO.PollStopButton(drvStop)

	go fsm.UpdateElevatorInformation(elevatorToNOCh, elevatorToFACh, elevatorToOECh, elevatorUpdateCh)
	go fsm.NewOrderFSM(buttonCh, elevatorToNOCh, elevatorUpdateCh, runTimerCh)
	go fsm.FloorArrivalFSM(floorCh, elevatorToFACh, elevatorUpdateCh, runTimerCh)
	go fsm.OrderExecutedFSM(timedOutCh, elevatorToOECh, elevatorUpdateCh)

	go timer.RunTimer(runTimerCh, timedOutCh)

	fsm.InitializeElevator(floorCh, elevatorUpdateCh)


	for {
		select {
		case a := <- drvButtons:
			fmt.Printf("%+v\n", a)
			hardwareIO.SetButtonLamp(a.Button, a.Floor, true)
			buttonCh <- a
		case a := <- drvFloors:
			fmt.Printf("%+v\n", a)
			floorCh <- a
			//if a == numFloors-1 {
			//	hardwareIO.SetMotorDirection(hardwareIO.MD_Down)
			//} else if a == 0 {
			//	hardwareIO.SetMotorDirection(hardwareIO.MD_Up)
			//}
		case a := <- drvObstr:
			fmt.Printf("%+v\n", a)
			if a {
				hardwareIO.SetMotorDirection(hardwareIO.MD_Stop)
			} else {
			}

		case a := <- drvStop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := hardwareIO.ButtonType(0); b < 3; b++ {
					hardwareIO.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}