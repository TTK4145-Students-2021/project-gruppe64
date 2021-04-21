package hardwareIO

import (
	"fmt"
	"realtimeProject/project-gruppe64/system"
)

/*
import (
	"../system"
	"fmt"
)
 */


func RunHardware(orderToSelfCh chan<- system.ButtonEvent, hallOrderCh chan<- system.ButtonEvent,
	floorArrivalCh chan<- int, obstructionEventCh chan<- bool)  {

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
			if a.Button == system.BT_Cab { //Sjekker om til fsm eller til distributor
				orderToSelfCh <- a
			} else {
				hallOrderCh <- a
			}
		case a := <- drvFloors:
			fmt.Printf("%+v\n", a)
			floorArrivalCh <- a
		case a := <- drvObstr:
			fmt.Printf("%+v\n", a)
			obstructionEventCh <- a
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
		default:
			break
		}
	}
}