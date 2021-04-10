package fsm

import (
	"fmt"
	"realtimeProject/project-gruppe64/configuration"
	"realtimeProject/project-gruppe64/hardwareIO"
)

func setAllButtonLights(e Elevator){
	for f := 0; f < configuration.NumFloors; f++ {
		for b := 0; b < configuration.NumButtons; b++  {
			if e.Orders[f][b] != 0 {
				hardwareIO.SetButtonLamp(hardwareIO.ButtonType(b), f, true)
			} else {
				hardwareIO.SetButtonLamp(hardwareIO.ButtonType(b), f, false)
			}
		}
	}
}

func orderAbove(e Elevator) bool {
	for f := e.Floor+1; f < configuration.NumFloors ; f++ {
		for b := 0; b < configuration.NumButtons; b++ {
			if e.Orders[f][b] != 0 {
				return true
			}
		}
	}
	return false
}

func orderBelow(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for b := 0; b < configuration.NumButtons; b++ {
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
		for button := 0; button < configuration.NumButtons; button++ { //_numButtons= 3
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


/////////////////////////////////////////////////////////////////////////////////////
////////////////////////////Maybe not necessary/////////////////////////////////////

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
	for f := configuration.NumFloors - 1; f >= 0; f-- {
		fmt.Printf("\n  | %d", f)
		for b := 0; b < configuration.NumButtons; b++ {
			if (f == configuration.NumFloors-1 && b == int(hardwareIO.BT_HallUp)) || (f == 0 && b == hardwareIO.BT_HallDown) {
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