package fsm


import (
	"../hardwareIO"
	"../system"
	"fmt"
)

// GOROUTINE, main initiated
// Controls the state of the elevator based on input from hardwareIO. Updates the designator on its own
// elevator struct. Uses the doorTimer for timing the opening of doors.
func ElevatorFSM(orderToSelfCh <-chan system.ButtonEvent, floorArrivalCh <-chan int, obstructionEventCh <-chan bool,
	ownElevatorCh chan<- system.Elevator, doorTimerDurationCh chan<- float64, doorTimerTimedOutCh <-chan bool,
	motorErrorCh <-chan bool, removeOrderCh <-chan system.ButtonEvent){
	var elevator system.Elevator
	obstruction := false
	elevator.ID = system.ElevatorID
	elevator.Orders = system.GetLoggedElevator().Orders

	select {
	case floorArrival :=<- floorArrivalCh:
		elevator.Floor = floorArrival
		elevator.MotorDirection = system.MDStop
		elevator.Behaviour = system.EBIdle
		elevator.Config.ClearOrdersVariant = system.ElevatorClearOrdersVariant
		elevator.Config.DoorOpenDurationSec = system.ElevatorDoorOpenDuration
		break
	default: // If no floor is detected by the floor sensor at start up
		elevator.Floor = -1
		elevator.MotorDirection = system.MDDown
		hardwareIO.SetMotorDirection(elevator.MotorDirection)
		elevator.Behaviour = system.EBMoving
		elevator.Config.ClearOrdersVariant = system.ElevatorClearOrdersVariant
		elevator.Config.DoorOpenDurationSec = system.ElevatorDoorOpenDuration
		break
	}

	for{
		select {
		case motorError := <-motorErrorCh:
			elevator.MotorError = motorError
			ownElevatorCh <- elevator

		case removeOrder := <-removeOrderCh:
			elevator.Orders[removeOrder.Floor][removeOrder.Button] = 0

		case orderToSelf := <-orderToSelfCh:
			hardwareIO.SetButtonLamp(orderToSelf.Button, orderToSelf.Floor, true)
			if obstruction{
				elevator.Orders[orderToSelf.Floor][int(orderToSelf.Button)] = 1
				break
			}
			switch elevator.Behaviour {
			case system.EBDoorOpen:
				if elevator.Floor == orderToSelf.Floor {
					doorTimerDurationCh <- elevator.Config.DoorOpenDurationSec
				} else {
					elevator.Orders[orderToSelf.Floor][int(orderToSelf.Button)] = 1
				}
				break

			case system.EBMoving:
				elevator.Orders[orderToSelf.Floor][int(orderToSelf.Button)] = 1
				break
			case system.EBIdle:
				if elevator.Floor == orderToSelf.Floor {
					hardwareIO.SetDoorOpenLamp(true)
					doorTimerDurationCh <- elevator.Config.DoorOpenDurationSec
					elevator.Behaviour = system.EBDoorOpen
					hardwareIO.SetMotorDirection(system.MDStop)
				} else {
					elevator.Orders[orderToSelf.Floor][int(orderToSelf.Button)] = 1
					elevator.MotorDirection = chooseDirection(elevator)
					hardwareIO.SetMotorDirection(elevator.MotorDirection)
					elevator.Behaviour = system.EBMoving
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
			case system.EBMoving:
				if elevatorShouldStop(elevator){
					hardwareIO.SetMotorDirection(system.MDStop)
					hardwareIO.SetDoorOpenLamp(true)
					elevator = clearOrdersAtCurrentFloor(elevator)
					doorTimerDurationCh <- elevator.Config.DoorOpenDurationSec
					setAllButtonLights(elevator)
					elevator.Behaviour = system.EBDoorOpen
				} else if elevator.Floor == 0{
					elevator.MotorDirection = system.MDUp
					hardwareIO.SetMotorDirection(elevator.MotorDirection)
				} else if elevator.Floor == 3 {
					elevator.MotorDirection = system.MDDown
					hardwareIO.SetMotorDirection(elevator.MotorDirection)
				} else if obstruction{
					elevator.MotorDirection = system.MDStop
					hardwareIO.SetMotorDirection(elevator.MotorDirection)
					hardwareIO.SetDoorOpenLamp(true)
					elevator.Behaviour = system.EBDoorOpen
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
				case system.EBDoorOpen:
					clearOrdersAtCurrentFloor(elevator)
					elevator.MotorDirection = chooseDirection(elevator)
					hardwareIO.SetDoorOpenLamp(false)
					hardwareIO.SetMotorDirection(elevator.MotorDirection)
					if elevator.MotorDirection == system.MDStop {
						elevator.Behaviour = system.EBIdle
					} else {
						elevator.Behaviour = system.EBMoving
					}
					break
				default:
					break
				}
			}
			ownElevatorCh <- elevator

		case obstructionEvent := <-obstructionEventCh:
			obstruction = obstructionEvent
			if elevator.Behaviour == system.MDStop && obstructionEvent{
				hardwareIO.SetDoorOpenLamp(true)
				elevator.Behaviour = system.EBDoorOpen
			}
			if !obstruction {
				doorTimerDurationCh <- elevator.Config.DoorOpenDurationSec
				break
			}
		}
	}
}
