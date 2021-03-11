package main

import (
	"realtimeProject/project-gruppe64/elevator"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/io"
)
import "fmt"

func main(){

	numFloors := elevator.NumFloors
	elev := elevator.ElevatorUnitialized()

	io.Init("localhost:15657", numFloors)
	fsm.Init(elev)

	var d io.MotorDirection = io.MD_Up
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan io.ButtonEvent)
	drv_floors  := make(chan int)
	drv_obstr   := make(chan bool)
	drv_stop    := make(chan bool)

	go io.PollButtons(drv_buttons)
	go io.PollFloorSensor(drv_floors)
	go io.PollObstructionSwitch(drv_obstr)
	go io.PollStopButton(drv_stop)


	for {
		select {
		case a := <- drv_buttons:
			fmt.Printf("%+v\n", a)
			io.SetButtonLamp(a.Button, a.Floor, true)
			fsm.FSMOnRequestButtonPress(a.Floor, a.Button)

		case a := <- drv_floors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				d = io.MD_Down
			} else if a == 0 {
				d = io.MD_Up
			}
			io.SetMotorDirection(d)
			fsm.FSMOnFloorArrival(a)


		case a := <- drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				io.SetMotorDirection(io.MD_Stop)
			} else {
				io.SetMotorDirection(d)
			}

		case a := <- drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := io.ButtonType(0); b < 3; b++ {
					io.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}