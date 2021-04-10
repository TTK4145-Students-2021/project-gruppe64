package hardwareIO

import (
	"fmt"
	"realtimeProject/project-gruppe64/system"
)




func RunHardware(orderToSelf chan<- system.ButtonEvent, hallOrder chan<- system.ButtonEvent, floorArrival chan<- int, obstructionEvent chan<- bool)  {

	drvButtons := make(chan system.ButtonEvent)
	drvFloors  := make(chan int)
	drvObstr   := make(chan bool)
	drvStop    := make(chan bool)

	go PollButtons(drvButtons)
	go PollFloorSensor(drvFloors)
	go PollObstructionSwitch(drvObstr)
	go PollStopButton(drvStop)

	for {
		select {
		case a := <- drvButtons:
			fmt.Printf("%+v\n", a)
			SetButtonLamp(a.Button, a.Floor, true)
			if a.Button == system.BT_Cab { //Sjekker om til fsm eller til distributor
				orderToSelf <- a
			} else {
				hallOrder <- a
			}
		case a := <- drvFloors:
			fmt.Printf("%+v\n", a)
			floorArrival <- a
		case a := <- drvObstr:
			fmt.Printf("%+v\n", a)
			obstructionEvent <- a
		case a := <- drvStop:
			// Can choose if implemented
			for a {
				SetMotorDirection(system.MD_Stop)
				SetStopLamp(true)
			}
			fmt.Printf("%+v\n", a)
			for f := 0; f < system.NumFloors; f++ {
				for b := system.ButtonType(0); b < 3; b++ {
					SetButtonLamp(b, f, false)
				}
			}
		}
	}
}