package hardwareIO

import (
	"fmt"
)

const (
	NumFloors = 4
	NumButtons = 3
)


func RunHardware(orderToSelf chan<- ButtonEvent, hallOrder chan<- ButtonEvent, floorArrival chan<- int)  {

	drvButtons := make(chan ButtonEvent)
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
			if a.Button == BT_Cab { //Sjekker om til fsm eller til distributor
				orderToSelf <- a
			} else {
				hallOrder <- a
			}
		case a := <- drvFloors:
			fmt.Printf("%+v\n", a)
			floorArrival <- a
		case a := <- drvObstr:
			fmt.Printf("%+v\n", a)
			//if a {
			//	SetMotorDirection(MD_Stop)
			//} else {}
		case a := <- drvStop:
			if a {
				SetMotorDirection(MD_Stop)
			} else {}
			fmt.Printf("%+v\n", a)
			for f := 0; f < NumFloors; f++ {
				for b := ButtonType(0); b < 3; b++ {
					SetButtonLamp(b, f, false)
				}
			}
		}
	}
}