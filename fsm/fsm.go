package fsm

import (
	"fmt"
	"realtimeProject/project-gruppe64/elevator"
	"realtimeProject/project-gruppe64/io"
	"realtimeProject/project-gruppe64/requests"
	"realtimeProject/project-gruppe64/timer"
)

var elev elevator.Elevator

func Init(e elevator.Elevator) {
	elev = e
}

func setAllLights(){ //Denne ble nok feil (?)
	for floor := 0; floor < elevator.NumFloors; floor++ {
		for btn := 0; btn < elevator.NumButtons; btn++ {
			io.SetButtonLamp(io.BT_Cab, floor, true ) //NOt sure if right button
		}
	}
}

func FSMOnInitBetweenFloors(){
	io.SetMotorDirection(io.MD_Down)
	elev.MotorDirection = io.MD_Down
	elev.Behaviour = elevator.EB_Moving
}

func FSMOnRequestButtonPress(btnFloor int, btnType io.ButtonType){
	//PRINT BUTTON FLOOR AND TYPE, use printf
	//PRINT ELEVATOR
	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		if elev.Floor == btnFloor {
			timer.TimerStart(elev.Config.DoorOpenDurationSec)
		} else {
			elev.Requests[btnFloor][btnType] = 1;
		}
		break

	case elevator.EB_Moving:
		elev.Requests[btnFloor][btnType] = 1;
		break

	case elevator.EB_Idle:
		if elev.Floor == btnFloor {
			io.SetDoorOpenLamp(true)
			timer.TimerStart(elev.Config.DoorOpenDurationSec)
			elev.Behaviour = elevator.EB_DoorOpen
		} else {
			elev.Requests[btnFloor][btnType] = 1
			elev.MotorDirection = requests.RequestsChooseDirection(elev)
			io.SetMotorDirection(elev.MotorDirection)
		}
		break
	}
	setAllLights()
	fmt.Printf("\nNew state:\n")
	elevator.ElevatorPrint(elev)
}

func FSMOnFloorArrival(newFloor int){
	//PRINT NEW FLOOR
	fmt.Printf("\n New floor \n")
	fmt.Printf("  |floor = %-2d          |\n",newFloor)
	elevator.ElevatorPrint(elev)
	elev.Floor = newFloor

	io.SetFloorIndicator(elev.Floor)
	//ElevoutputDevice.floorIndicator(elev.floor)

	switch elev.Behaviour {
	case elevator.EB_Moving:
		if requests.RequestsShouldStop(elev){
			io.SetMotorDirection(io.MD_Stop)
			io.SetDoorOpenLamp(true)
			elev = requests.RequestsClearAtCurrentFloor(elev)
			timer.TimerStart(elev.Config.DoorOpenDurationSec)
			setAllLights()
			elev.Behaviour = elevator.EB_DoorOpen
		}
	default:
	}
	fmt.Println("\nNew state:")
	elevator.ElevatorPrint(elev)
}

func FSMOnDoorTimeout(){
	// PRINT SOME FORMATTING STUFF
	elevator.ElevatorPrint(elev)

	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		elev.MotorDirection = requests.RequestsChooseDirection(elev)

		io.SetDoorOpenLamp(false)
		io.SetMotorDirection(elev.MotorDirection)

		if elev.MotorDirection == io.MD_Stop {
			elev.Behaviour = elevator.EB_Idle
		} else {
			elev.Behaviour = elevator.EB_Moving
		}

		break
	default:
		break
	}

	fmt.Printf("\nNew state:\n")
	elevator.ElevatorPrint(elev)
}
