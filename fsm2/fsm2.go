package fsm2

import (
	"fmt"
	"../hardwareIO"
)

type ElevatorBehaviour int
const (
	EB_Idle     ElevatorBehaviour = 0 // Evt skrive om til camelCase!
	EB_DoorOpen                   = 1
	EB_Moving                     = 2
)

type ClearOrdersVariant int
const (
	CO_All    ClearOrdersVariant = 0
	CO_InMotorDirection                     = 1
)

type Elevator struct {
	Floor          int
	MotorDirection hardwareIO.MotorDirection
	Orders       [hardwareIO.NumFloors][hardwareIO.NumButtons] int
	Behaviour      ElevatorBehaviour
	Config         struct{
		ClearOrdersVariant ClearOrdersVariant
		DoorOpenDurationSec float64
	}
}

func elevatorBehaviourToString(eB ElevatorBehaviour) string {
	switch eB {
	case EB_Idle:
		return "ED_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	}
	return "EB_Undefined"
}

func motorDirectionToString(mD hardwareIO.MotorDirection) string { // Importere? Stor forbokstav!
	switch mD {
	case hardwareIO.MD_Stop:
		return "MD_Stop"
	case hardwareIO.MD_Up:
		return "MD_Up"
	case hardwareIO.MD_Down:
		return "MD_Down"
	}
	return "MD_Undefined"
}

func printElevator(e Elevator) {
	fmt.Printf("  +--------------------+\n")                                             // Sjekk om dette er riktig print-funksjon!
	fmt.Printf("  |floor = %-2d          |\n", e.Floor)                               // - i %2-d betyr bare - at teksten er left-justified (kosmetisk)
	fmt.Printf("|direction  = %-12.12s|\n", motorDirectionToString(e.MotorDirection)) // Måtte dele opp for at det skulle bli riktig
	fmt.Printf("|behaviour = %s|\n", elevatorBehaviourToString(e.Behaviour))                         // Hvorfor feilmelding her?
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("  |  | up  | dn  | cab |\n")
	for f := hardwareIO.NumFloors - 1; f >= 0; f-- {
		fmt.Printf("\n  | %d", f)
		for b := 0; b < hardwareIO.NumButtons; b++ {
			if (f == hardwareIO.NumFloors-1 && b == int(hardwareIO.BT_HallUp)) || (f == 0 && b == hardwareIO.BT_HallDown) {
				fmt.Printf("|     ")
			} else {
				if e.Orders[f][b] != 0 {
					fmt.Printf("|  #  ")
				} else {
					fmt.Printf("|  -  ")
				}
			}
		}
	}
}

func orderAbove(e Elevator) bool {
	for f := e.Floor+1; f < hardwareIO.NumFloors ; f++ {
		for b := 0; b < hardwareIO.NumButtons; b++ {
			if e.Orders[f][b] != 0 {
				return true
			}
		}
	}
	return false
}

func orderBelow(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for b := 0; b < hardwareIO.NumButtons; b++ {
			if e.Orders[f][b] != 0 {
				return true
			}
		}
	}
	return false
}

func chooseDirection(e Elevator) hardwareIO.MotorDirection{
	var returnDir hardwareIO.MotorDirection
	switch e.MotorDirection {
	case hardwareIO.MD_Up:
		if orderAbove(e){
			returnDir = hardwareIO.MD_Up
		} else if orderBelow(e) {
			returnDir = hardwareIO.MD_Down
		} else {
			returnDir = hardwareIO.MD_Stop
		}
	case hardwareIO.MD_Down:

		if orderBelow(e){
			returnDir = hardwareIO.MD_Down
		} else if orderAbove(e){
			returnDir = hardwareIO.MD_Up
		} else {
			returnDir = hardwareIO.MD_Stop
		}
	case hardwareIO.MD_Stop:
		if orderBelow(e){
			returnDir = hardwareIO.MD_Down
		} else if orderAbove(e){
			returnDir = hardwareIO.MD_Up
		}
	default:
		break
	}
	return returnDir
}

func elevatorShouldStop(e Elevator) bool {
	switch e.MotorDirection {
	case hardwareIO.MD_Down:
		if e.Orders[e.Floor][hardwareIO.BT_HallDown] != 0 || e.Orders[e.Floor][hardwareIO.BT_Cab] != 0 || !orderBelow(e){ //This gives me if eRequests == true, right?
			return true
		} else{
			return false
		}
	case hardwareIO.MD_Up:
		if e.Orders[e.Floor][hardwareIO.BT_HallUp] != 0 || e.Orders[e.Floor][hardwareIO.BT_Cab] != 0 || !orderAbove(e){ //This gives me if eRequests == true, right?
			return true
		} else{
			return false
		}

	case hardwareIO.MD_Stop:
		break
	default:
		break
	}
	return true
}

func clearOrdersAtCurrentFloor(e Elevator) Elevator{
	switch e.Config.ClearOrdersVariant {
	case CO_All: //CV:clear request variant
		for button := 0; button < hardwareIO.NumButtons; button++ { //_numButtons= 3
			e.Orders[e.Floor][button] = 0
		}

	case CO_InMotorDirection:
		e.Orders[e.Floor][hardwareIO.BT_Cab] = 0
		switch e.MotorDirection {
		case hardwareIO.MD_Up:
			e.Orders[e.Floor][hardwareIO.BT_HallUp] = 0
			if orderAbove(e) == false {
				e.Orders[e.Floor][hardwareIO.BT_HallDown] = 0
			}

		case hardwareIO.MD_Down:
			e.Orders[e.Floor][hardwareIO.BT_HallDown] = 0
			if orderBelow(e) == false {
				e.Orders[e.Floor][hardwareIO.BT_HallUp] = 0
			}

		case hardwareIO.MD_Stop:
		default:
			e.Orders[e.Floor][hardwareIO.BT_HallUp] = 0
			e.Orders[e.Floor][hardwareIO.BT_HallDown] = 0
		}
	}
	return e
}

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