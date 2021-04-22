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
	"time"
)
 */

// GO-ROUTINE, main initiated
// Based on https://github.com/TTK4145/driver-go/blob/master/main.go. Hall button-orders are sent to
// distributor and cab button-orders are sent to FSM
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
			if a.Button == system.BTCab {
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

// GO-ROUTINE, main initiated
// Checks the elevator log for motor stop every system-given time.
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
					if !elevatorNow.MotorError{
						motorErrorCh <- true
						fmt.Println("Motor error")
					}
				}
			} else {
				if elevatorNow.MotorError{
					motorErrorCh <- false
					fmt.Println("Motor functioning again!")
				}
			}
		})
	}
}