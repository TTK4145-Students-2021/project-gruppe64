package fsm

import (
	"fmt"
	"realtimeProject/project-gruppe64/io"
	"realtimeProject/project-gruppe64/timer"
)

var elevator Elevator

func setAllLights(es Elevator){ //Denne ble nok feil (?)
	for floor := 0; floor < NFloors; floor++ {
		for btn := 0; btn < NButtons; btn++ {
			io.SetButtonLamp(io.BT_Cab, floor, true ) //NOt sure if right button
		}
	}
}

func FSMOnInitBetweenFloors(){
	io.SetMotorDirection(io.MD_Down)
	elevator.Dirn = io.MD_Down
	elevator.Behaviour = EB_Moving
}

func FSMOnRequestButtonPress(btnFloor int, btnType io.ButtonType){
	//PRINT BUTTON FLOOR AND TYPE, use printf
	//PRINT ELEVATOR
	switch elevator.Behaviour {
	case EB_DoorOpen:
		if elevator.Floor == btnFloor {
			timer.TimerStart(elevator.Config.DoorOpenDuration_s)
		} else {
			elevator.Requests[btnFloor][btnType] = 1;
		}
		break

	case EB_Moving:
		elevator.Requests[btnFloor][btnType] = 1;
		break

	case EB_Idle:
		if elevator.Floor == btnFloor {
			io.SetDoorOpenLamp(true)
			timer.TimerStart(elevator.Config.DoorOpenDuration_s)
			elevator.Behaviour = EB_DoorOpen
		} else {
			elevator.Requests[btnFloor][btnType] = 1
			elevator.Dirn = RequestsChooseDirection(elevator)
			io.SetMotorDirection(elevator.Dirn)
		}
		break
	}
	setAllLights(elevator)
	fmt.Printf("\nNew state:\n")
	ElevatorPrint(elevator)
}

func FSMOnFloorArrival(newFloor int){
	//PRINT NEW FLOOR
	ElevatorPrint(elevator)
	elevator.Floor = newFloor

	io.SetFloorIndicator(elevator.Floor)
	//ElevoutputDevice.floorIndicator(elevator.floor)

	switch elevator.behaviour {
	case EB_Moving:
		if requests.RequestsShouldStop(elevator){
			io.SetMotorDirection(io.MD_Stop)
			io.SetDoorOpenLamp(true)
			elevator = RequestsClearAtCurrentFloor(elevator)
			timer.TimerStart(elevator.Config.DoorOpenDuration_s)
			setAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen
		}
	default:
	}
	fmt.Println("\nNew state:")
	ElevatorPrint(elevator)
}

func FSMOnDoorTimeout(){
	// PRINT SOME FORMATTING STUFF
	ElevatorPrint(elevator)

	switch elevator.Behaviour {
	case EB_DoorOpen:
		elevator.Dirn = RequestsChooseDirection(elevator)

		io.SetDoorOpenLamp(false)
		io.SetMotorDirection(elevator.Dirn)

		if elevator.Dirn == io.MD_Stop {
			elevator.Behaviour = EB_Idle
		} else {
			elevator.Behaviour = EB_Moving
		}

		break
	default:
		break
	}

	fmt.Printf("\nNew state:\n")
	ElevatorPrint(elevator)
}

