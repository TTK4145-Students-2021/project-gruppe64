package hardwareIO

import (
	"fmt"
	"realtimeProject/project-gruppe64/system"
	"time"
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
			if a.Button == system.BTCab { //Sjekker om til fsm eller til distributor
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
			fmt.Printf("%+v\n", a)
		default:
			break
		}
	}
}

func CheckForMotorStop(motorErrorCh chan <- bool){
	for {
		time.Sleep(time.Millisecond * 500)
		elevatorBefore := system.GetLoggedElevator()
		time.AfterFunc(system.CheckMotorAfterDuration*time.Second, func() {
			elevatorNow := system.GetLoggedElevator()
			if elevatorBefore.Floor == elevatorNow.Floor && elevatorBefore.Behaviour == elevatorNow.Behaviour &&
				elevatorBefore.MotorDirection == elevatorNow.MotorDirection {
				ordersBeforeNum := 0
				ordersNowNum := 0
				for f := 0; f < system.NumFloors; f++ {
					for b := 0; b < system.NumButtons; b++ {
						ordersBeforeNum += elevatorBefore.Orders[f][b]
						ordersNowNum += elevatorNow.Orders[f][b]
					}
				}
				if ordersBeforeNum >= ordersNowNum && ordersBeforeNum > 0{
					motorErrorCh <- true
					fmt.Println("Motor error")
				} else {
					motorErrorCh <- false
				}
			} else {
				motorErrorCh <- false
			}
		})
	}
}