package fsm

import (
	"fmt"
	"realtimeProject/project-gruppe64/elevator"
	"realtimeProject/project-gruppe64/io"
	"realtimeProject/project-gruppe64/requests"
	"realtimeProject/project-gruppe64/timer"
)

var runningElevator elevator.Elevator

func Init() {
	runningElevator = elevator.ElevatorUnitialized()
}

func FSMOnInitBetweenFloors(){
	io.SetMotorDirection(io.MD_Down)
	runningElevator.MotorDirection = io.MD_Down
	runningElevator.Behaviour = elevator.EB_Moving
}

func setAllLights(){
	for f := 0; f < elevator.NumFloors; f++{
		for b := 0; b < elevator.NumButtons; b++ {
			if runningElevator.Requests[f][b] != 0 {
				io.SetButtonLamp(io.ButtonType(b), f, true)
			} else {
				io.SetButtonLamp(io.ButtonType(b), f, true)

			}
		}
	}
}

func ElevioButtonToString(btnType io.ButtonType) string {
	switch btnType {
	case io.BT_HallUp:
		return "BT_HallUp"
	case io.BT_HallDown:
		return "BT_HallDown"
	case io.BT_Cab:
		return "BT_Cab"
	}
	return "BT_Undefined"
}

func FSMOnRequestButtonPress(btnFloor int, btnType io.ButtonType){
	fmt.Printf("\n \n %-2d", btnFloor)
	fmt.Printf(" %-12.12s \n", ElevioButtonToString(btnType))

	switch runningElevator.Behaviour {
	case elevator.EB_DoorOpen:
		if runningElevator.Floor == btnFloor {
			timer.TimerStart(runningElevator.Config.DoorOpenDurationSec)
		} else {
			runningElevator.Requests[btnFloor][int(btnType)] = 1
		}
		break

	case elevator.EB_Moving:
		runningElevator.Requests[btnFloor][int(btnType)] = 1
		break

	case elevator.EB_Idle:
		if runningElevator.Floor == btnFloor {
			io.SetDoorOpenLamp(true)
			timer.TimerStart(runningElevator.Config.DoorOpenDurationSec)
			runningElevator.Behaviour = elevator.EB_DoorOpen
		} else {
			runningElevator.Requests[btnFloor][int(btnType)] = 1
			runningElevator.MotorDirection = requests.RequestsChooseDirection(runningElevator)
			io.SetMotorDirection(runningElevator.MotorDirection)
		}
		break
	}
	setAllLights()
	fmt.Printf("\nNew state:\n")
	elevator.ElevatorPrint(runningElevator)
}

func FSMOnFloorArrival(newFloor int){
	//PRINT NEW FLOOR
	fmt.Printf("\n \n New floor: %-2d \n",newFloor)
	elevator.ElevatorPrint(runningElevator)

	runningElevator.Floor = newFloor

	io.SetFloorIndicator(runningElevator.Floor)

	switch runningElevator.Behaviour {
	case elevator.EB_Moving:
		if requests.RequestsShouldStop(runningElevator){
			io.SetMotorDirection(io.MD_Stop)
			io.SetDoorOpenLamp(true)
			runningElevator = requests.RequestsClearAtCurrentFloor(runningElevator)
			timer.TimerStart(runningElevator.Config.DoorOpenDurationSec)
			setAllLights()
			runningElevator.Behaviour = elevator.EB_DoorOpen
		}
		break
	default:
		break
	}
	fmt.Println("\nNew state:")
	elevator.ElevatorPrint(runningElevator)
}

func FSMOnDoorTimeout(){
	elevator.ElevatorPrint(runningElevator)

	switch runningElevator.Behaviour {
	case elevator.EB_DoorOpen:
		runningElevator.MotorDirection = requests.RequestsChooseDirection(runningElevator)

		io.SetDoorOpenLamp(false)
		io.SetMotorDirection(runningElevator.MotorDirection)

		if runningElevator.MotorDirection == io.MD_Stop {
			runningElevator.Behaviour = elevator.EB_Idle
		} else {
			runningElevator.Behaviour = elevator.EB_Moving
		}

		break
	default:
		break
	}

	fmt.Printf("\nNew state:\n")
	elevator.ElevatorPrint(runningElevator)
}