package fsm

import (
	"fmt"
	"realtimeProject/project-gruppe64/system"
	"realtimeProject/project-gruppe64/hardwareIO"
)

/*
import (
	"../hardwareIO"
	"../system"
	"fmt"
)
 */


func ElevatorFSM(orderToSelf <-chan system.ButtonEvent, floorArrival <-chan int, obstructionEvent <-chan bool, ownElevator chan<- system.Elevator, doorTimerDuration chan<- float64, doorTimerTimedOut <-chan bool){
	var elevator system.Elevator
	obstruction := false
	elevator.ID = system.ElevatorID
	elevator.Orders = system.GetLoggedOrders()
	select {
	case flrA :=<- floorArrival: // If the floor sensor registers a floor at initialization
		elevator.Floor = flrA
		elevator.MotorDirection = system.MD_Stop
		elevator.Behaviour = system.EB_Idle
		elevator.Config.ClearOrdersVariant = system.CO_InMotorDirection
		elevator.Config.DoorOpenDurationSec = 3.0
		break
	default: // If no floor is detected by the floor sensor
		elevator.Floor = -1
		elevator.MotorDirection = system.MD_Down
		hardwareIO.SetMotorDirection(system.MD_Down)
		elevator.Behaviour = system.EB_Moving
		elevator.Config.ClearOrdersVariant = system.CO_InMotorDirection
		elevator.Config.DoorOpenDurationSec = 3.0
		break
	}

	for{
		select {
		case btnE := <-orderToSelf:
			hardwareIO.SetButtonLamp(btnE.Button, btnE.Floor, true)
			if obstruction{
				elevator.Orders[btnE.Floor][int(btnE.Button)] = 1
				break
			}
			switch elevator.Behaviour {
			case system.EB_DoorOpen:
				if elevator.Floor == btnE.Floor {
					doorTimerDuration <- elevator.Config.DoorOpenDurationSec
				} else {
					elevator.Orders[btnE.Floor][int(btnE.Button)] = 1
				}
				break

			case system.EB_Moving:
				elevator.Orders[btnE.Floor][int(btnE.Button)] = 1
				break
			case system.EB_Idle:
				if elevator.Floor == btnE.Floor {
					hardwareIO.SetDoorOpenLamp(true)
					doorTimerDuration <- elevator.Config.DoorOpenDurationSec
					elevator.Behaviour = system.EB_DoorOpen
					hardwareIO.SetMotorDirection(system.MD_Stop)
				} else {
					elevator.Orders[btnE.Floor][int(btnE.Button)] = 1
					elevator.MotorDirection = chooseDirection(elevator)
					hardwareIO.SetMotorDirection(elevator.MotorDirection)
					elevator.Behaviour = system.EB_Moving
				}
				break
			default:
				fmt.Printf("\n Button was bushed but nothing happend. Undefined state.\n")
				break
			}
			ownElevator <- elevator
		case flrA := <-floorArrival:
			elevator.Floor = flrA
			hardwareIO.SetFloorIndicator(elevator.Floor)
			switch elevator.Behaviour {
			case system.EB_Moving:
				if elevatorShouldStop(elevator){
					hardwareIO.SetMotorDirection(system.MD_Stop)
					hardwareIO.SetDoorOpenLamp(true)
					elevator = clearOrdersAtCurrentFloor(elevator)
					doorTimerDuration <- elevator.Config.DoorOpenDurationSec
					setAllButtonLights(elevator)
					elevator.Behaviour = system.EB_DoorOpen
				} else if elevator.Floor == 0{
					elevator.MotorDirection = system.MD_Up
					hardwareIO.SetMotorDirection(system.MD_Up)
				} else if elevator.Floor == 3 {
					elevator.MotorDirection = system.MD_Down
					hardwareIO.SetMotorDirection(system.MD_Down)
				} else if obstruction{
					hardwareIO.SetMotorDirection(system.MD_Stop)
					hardwareIO.SetDoorOpenLamp(true)
					elevator.Behaviour = system.EB_DoorOpen
				}
				break
			default:
				fmt.Printf("\n Arrived at floor but nothing happend. Undefined state.\n")
				break
			}
			setAllButtonLights(elevator)
			ownElevator <- elevator
		case dTTimedOut := <-doorTimerTimedOut:
			if obstruction{
				break
			}
			if dTTimedOut {
				switch elevator.Behaviour {
				case system.EB_DoorOpen:
					clearOrdersAtCurrentFloor(elevator)
					elevator.MotorDirection = chooseDirection(elevator)
					hardwareIO.SetDoorOpenLamp(false)
					hardwareIO.SetMotorDirection(elevator.MotorDirection)
					if elevator.MotorDirection == system.MD_Stop {
						elevator.Behaviour = system.EB_Idle
					} else {
						elevator.Behaviour = system.EB_Moving
					}
					break
				default:
					break
				}
			}
			ownElevator <- elevator
		case obstrE := <-obstructionEvent:
			obstruction = obstrE
			if elevator.Behaviour == system.MD_Stop && obstrE{
				hardwareIO.SetDoorOpenLamp(true)
				elevator.Behaviour = system.EB_DoorOpen
			}

			if !obstruction {
				doorTimerDuration <- elevator.Config.DoorOpenDurationSec
				break
			}
		}
	}
}








