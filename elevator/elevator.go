package elevator

import (
	"fmt"
	"realtimeProject/project-gruppe64/io"
)
// Enum definitions

const (
	NumButtons int = 3 // Fra elevator_io_types.h
	NumFloors int = 4 // Skrevet om fra elevator_io - må kanskje endres
)

type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0 // Evt skrive om til camelCase!
	EB_DoorOpen                   = 1
	EB_Moving                     = 2
)

type ClearRequestVariant int

const (
	CV_All    ClearRequestVariant = 0
	CV_InDirn                     = 1
)

// Definition of structs

type Elevator struct {
	Floor          int
	MotorDirection io.MotorDirection
	Requests       [NumFloors][NumButtons] int // Hardkoder størrelse sånn halvveis
	Behaviour      ElevatorBehaviour
	Config         struct{
		ClearRequestVariant ClearRequestVariant
		DoorOpenDurationSec float64 // Tilsvarer double i C
	}
}

func ebToString(elevBehav ElevatorBehaviour) string {
	switch elevBehav {
	case EB_Idle:
		return "ED_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	}
	return "EB_Undefined"
}

func elevioMotorDirToString(mDirection io.MotorDirection) string { // Importere? Stor forbokstav!
	switch mDirection {
	case io.MD_Stop:
		return "MD_Stop"
	case io.MD_Up:
		return "MD_Up"
	case io.MD_Down:
		return "MD_Down"
	}
	return "MD_Undefined"
}


func ElevatorPrint(elev Elevator) {
	fmt.Printf("  +--------------------+\n")                                             // Sjekk om dette er riktig print-funksjon!
	fmt.Printf("  |floor = %-2d          |\n", elev.Floor)                               // - i %2-d betyr bare - at teksten er left-justified (kosmetisk)
	fmt.Printf("|direction  = %-12.12s|\n", elevioMotorDirToString(elev.MotorDirection)) // Måtte dele opp for at det skulle bli riktig
	fmt.Printf("|behaviour = %s|\n", ebToString(elev.Behaviour))                         // Hvorfor feilmelding her?
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("  |  | up  | dn  | cab |\n")
	for f := NumFloors - 1; f >= 0; f-- {
		fmt.Printf("\n  | %d", f)
		for b := 0; b < NumButtons; b++ {
			if (f == NumFloors-1 && b == int(io.BT_HallUp)) || (f == 0 && b == io.BT_HallDown) {
				fmt.Printf("|     ")
			} else {
				if elev.Requests[f][b] != 0 {
					fmt.Printf("|  #  ")
				} else {
					fmt.Printf("|  -  ")
				}
			}
		}
	}
}




func ElevatorUnitialized() Elevator { //
	elevator := Elevator{}
	elevator.Floor = -1
	elevator.MotorDirection = io.MD_Stop
	elevator.Behaviour = EB_Idle
	elevator.Config.ClearRequestVariant = CV_All
	elevator.Config.DoorOpenDurationSec = 3.0
	return elevator
}

