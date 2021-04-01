package fsm

import (
	"fmt"
	"realtimeProject/project-gruppe64/hardwareIO"
)



func ElevatorFSM(buttonEvent <-chan hardwareIO.ButtonEvent, floorArrival <-chan int, timerDuration chan<- float64, timedOut <-chan bool ){
	elevator := Elevator{}

	select {
	case flrA :=<- floorArrival: // If the floor sensor registers a floor at initialization
		elevator.Floor = flrA
		elevator.MotorDirection = hardwareIO.MD_Stop
		elevator.Behaviour = EB_Idle
		elevator.Config.ClearOrdersVariant = CO_InMotorDirection
		elevator.Config.DoorOpenDurationSec = 3.0
		break
	default: // If no floor is detected by the floor sensor
		elevator.Floor = -1
		elevator.MotorDirection = hardwareIO.MD_Down
		hardwareIO.SetMotorDirection(hardwareIO.MD_Down)
		elevator.Behaviour = EB_Moving
		elevator.Config.ClearOrdersVariant = CO_InMotorDirection
		elevator.Config.DoorOpenDurationSec = 3.0
		break
	}

	for{
		 // til en eller annen kanal til network goroutine for broadcasting av elevator struct
		select {
		case btnE := <-buttonEvent:
			hardwareIO.SetButtonLamp(btnE.Button, btnE.Floor, true)
			switch elevator.Behaviour{
			case EB_DoorOpen:
				if elevator.Floor == btnE.Floor {
					timerDuration <- elevator.Config.DoorOpenDurationSec
				} else {
					elevator.Orders[btnE.Floor][int(btnE.Button)] = 1
				}
				break
			case EB_Moving:
				elevator.Orders[btnE.Floor][int(btnE.Button)] = 1
				break
			case EB_Idle:
				if elevator.Floor == btnE.Floor {
					hardwareIO.SetDoorOpenLamp(true)
					timerDuration <- elevator.Config.DoorOpenDurationSec
					elevator.Behaviour = EB_DoorOpen
				} else {
					elevator.Orders[btnE.Floor][int(btnE.Button)] = 1
					elevator.MotorDirection = chooseDirection(elevator)
					hardwareIO.SetMotorDirection(elevator.MotorDirection)
					elevator.Behaviour = EB_Moving
				}
				break
			default:
				fmt.Printf("\n Button was bushed but nothing happend. Undefined state.\n")
				break
			}
		case flrA := <-floorArrival:
			elevator.Floor = flrA
			hardwareIO.SetFloorIndicator(elevator.Floor)
			switch elevator.Behaviour {
			case EB_Moving:
				if elevatorShouldStop(elevator){
					hardwareIO.SetMotorDirection(hardwareIO.MD_Stop)
					hardwareIO.SetDoorOpenLamp(true)
					elevator = clearOrdersAtCurrentFloor(elevator)
					timerDuration <- elevator.Config.DoorOpenDurationSec
					setAllButtonLights(elevator)
					elevator.Behaviour = EB_DoorOpen
				} else if elevator.Floor == 0{
					elevator.MotorDirection = hardwareIO.MD_Up
				} else if elevator.Floor == 3 {
					elevator.MotorDirection = hardwareIO.MD_Down
				}
				break
			default:
				fmt.Printf("\n Arrived at floor but nothing happend. Undefined state.\n")
				break
			}
			setAllButtonLights(elevator)
		case tmdO := <-timedOut:
			if tmdO {
				switch elevator.Behaviour {
				case EB_DoorOpen:
					clearOrdersAtCurrentFloor(elevator)
					elevator.MotorDirection = chooseDirection(elevator)
					hardwareIO.SetDoorOpenLamp(false)
					hardwareIO.SetMotorDirection(elevator.MotorDirection)
					if elevator.MotorDirection == hardwareIO.MD_Stop {
						elevator.Behaviour = EB_Idle
					} else {
						elevator.Behaviour = EB_Moving
					}
					break
				default:
					fmt.Printf("\n Timer timed out but nothing happend.:\n")
					break
				}
			}
		default:
			break
		}
	}
}


func setAllButtonLights(e Elevator){
	for f := 0; f < hardwareIO.NumFloors; f++ {
		for b := 0; b < hardwareIO.NumButtons; b++  {
			if e.Orders[f][b] != 0 {
				hardwareIO.SetButtonLamp(hardwareIO.ButtonType(b), f, true)
			} else {
				hardwareIO.SetButtonLamp(hardwareIO.ButtonType(b), f, false)
			}
		}
	}
}

