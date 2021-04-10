package fsm

import (
	"realtimeProject/project-gruppe64/configuration"
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
	Orders       [configuration.NumFloors][configuration.NumButtons] int
	Behaviour      ElevatorBehaviour
	Config         struct{
		ClearOrdersVariant ClearOrdersVariant
		DoorOpenDurationSec float64
	}
}

