package hardwareIO

import "fmt"

const (
	NumFloors = 4
	NumButtons = 3
)


func StartHardware(buttonPressed chan<- ButtonEvent, floorArrival chan<- int)  {
	Init("localhost:15657", NumFloors)

	drvButtons := make(chan ButtonEvent)
	drvFloors  := make(chan int)
	drvObstr   := make(chan bool)
	drvStop    := make(chan bool)

	go PollButtons(drvButtons)
	go PollFloorSensor(drvFloors)
	go PollObstructionSwitch(drvObstr)
	go PollStopButton(drvStop)

	go evaluateHardwareInputs(drvButtons, drvFloors, drvObstr, drvObstr, buttonPressed, floorArrival)
}

func evaluateHardwareInputs(drvB <-chan ButtonEvent, drvF <-chan int, drvO <-chan bool, drvS <-chan bool, btnP chan<- ButtonEvent, flrA chan<- int)  {
	for {
		select {
		case a := <- drvB:
			fmt.Printf("%+v\n", a)
			SetButtonLamp(a.Button, a.Floor, true)
			btnP <- a
		case a := <- drvF:
			fmt.Printf("%+v\n", a)
			flrA <- a
		case a := <- drvO:
			fmt.Printf("%+v\n", a)
			if a {
				SetMotorDirection(MD_Stop)
			} else {
			}
		case a := <- drvS:
			fmt.Printf("%+v\n", a)
			for f := 0; f < NumFloors; f++ {
				for b := ButtonType(0); b < 3; b++ {
					SetButtonLamp(b, f, false)
				}
			}
		}
	}
}