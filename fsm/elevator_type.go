package fsm

import (
	"fmt"
	"realtimeProject/project-gruppe64/hardwareIO"
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