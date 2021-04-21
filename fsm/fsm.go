package fsm

import (
	"fmt"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/system"
)

/*
import (
	"../hardwareIO"
	"../system"
	"fmt"
)
 */


func ElevatorFSM(orderToSelfCh <-chan system.ButtonEvent, floorArrivalCh <-chan int, obstructionEventCh <-chan bool,
	ownElevatorCh chan<- system.Elevator, doorTimerDurationCh chan<- float64, doorTimerTimedOutCh <-chan bool){

	var elevator system.Elevator
	obstruction := false
	elevator.ID = system.ElevatorID
	elevator.Orders = system.GetLoggedOrders()
	select {
	case floorArrival :=<- floorArrivalCh: // If the floor sensor registers a floor at initialization
		elevator.Floor = floorArrival
		elevator.MotorDirection = system.MD_Stop
		elevator.Behaviour = system.EB_Idle
		elevator.Config.ClearOrdersVariant = system.ElevatorClearOrdersVariant
		elevator.Config.DoorOpenDurationSec = system.ElevatorDoorOpenDuration
		break
	default: // If no floor is detected by the floor sensor
		elevator.Floor = -1
		elevator.MotorDirection = system.MD_Down
		hardwareIO.SetMotorDirection(system.MD_Down)
		elevator.Behaviour = system.EB_Moving
		elevator.Config.ClearOrdersVariant = system.ElevatorClearOrdersVariant
		elevator.Config.DoorOpenDurationSec = system.ElevatorDoorOpenDuration
		break
	}

	for{
		select {
		case orderToSelf := <-orderToSelfCh:
			hardwareIO.SetButtonLamp(orderToSelf.Button, orderToSelf.Floor, true)
			if obstruction{
				elevator.Orders[orderToSelf.Floor][int(orderToSelf.Button)] = 1
				break
			}
			switch elevator.Behaviour {
			case system.EB_DoorOpen:
				if elevator.Floor == orderToSelf.Floor {
					doorTimerDurationCh <- elevator.Config.DoorOpenDurationSec
				} else {
					elevator.Orders[orderToSelf.Floor][int(orderToSelf.Button)] = 1
				}
				break

			case system.EB_Moving:
				elevator.Orders[orderToSelf.Floor][int(orderToSelf.Button)] = 1
				break
			case system.EB_Idle:
				if elevator.Floor == orderToSelf.Floor {
					hardwareIO.SetDoorOpenLamp(true)
					doorTimerDurationCh <- elevator.Config.DoorOpenDurationSec
					elevator.Behaviour = system.EB_DoorOpen
					hardwareIO.SetMotorDirection(system.MD_Stop)
				} else {
					elevator.Orders[orderToSelf.Floor][int(orderToSelf.Button)] = 1
					elevator.MotorDirection = chooseDirection(elevator)
					hardwareIO.SetMotorDirection(elevator.MotorDirection)
					elevator.Behaviour = system.EB_Moving
				}
				break
			default:
				fmt.Printf("\n Button was bushed but nothing happend. Undefined state.\n")
				break
			}
			ownElevatorCh <- elevator
		case floorArrival := <-floorArrivalCh:
			elevator.Floor = floorArrival
			hardwareIO.SetFloorIndicator(elevator.Floor)
			switch elevator.Behaviour {
			case system.EB_Moving:
				if elevatorShouldStop(elevator){
					hardwareIO.SetMotorDirection(system.MD_Stop)
					hardwareIO.SetDoorOpenLamp(true)
					elevator = clearOrdersAtCurrentFloor(elevator)
					doorTimerDurationCh <- elevator.Config.DoorOpenDurationSec
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
			ownElevatorCh <- elevator
		case doorTimerTimedOut := <-doorTimerTimedOutCh:
			if obstruction{
				break
			}
			if doorTimerTimedOut {
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
			ownElevatorCh <- elevator
		case obstructionEvent := <-obstructionEventCh:
			obstruction = obstructionEvent
			if elevator.Behaviour == system.MD_Stop && obstructionEvent{
				hardwareIO.SetDoorOpenLamp(true)
				elevator.Behaviour = system.EB_DoorOpen
			}

			if !obstruction {
				doorTimerDurationCh <- elevator.Config.DoorOpenDurationSec
				break
			}
		}
	}
}








