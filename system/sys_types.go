package system

// Type based on type in https://github.com/TTK4145/driver-go/blob/master/elevio/elevator_io.go
type MotorDirection int
const (
	MDUp   MotorDirection = 1
	MDDown                = -1
	MDStop                = 0
)

// Type based on type in https://github.com/TTK4145/driver-go/blob/master/elevio/elevator_io.go
type ButtonType int
const (
	BTHallUp   ButtonType = 0
	BTHallDown            = 1
	BTCab                 = 2
)

// Type based on type in https://github.com/TTK4145/driver-go/blob/master/elevio/elevator_io.go
type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

// Type based on type in https://github.com/TTK4145/Project-resources/blob/master/elev_algo/elevator.h
type ElevatorBehaviour int
const (
	EBIdle     ElevatorBehaviour = 0
	EBDoorOpen                   = 1
	EBMoving                     = 2
)

// Type based on type in https://github.com/TTK4145/Project-resources/blob/master/elev_algo/elevator.h
type ClearOrdersVariant int
const (
	COAll    ClearOrdersVariant = 0
	COInMotorDirection          = 1
)

// Type based on type in https://github.com/TTK4145/Project-resources/blob/master/elev_algo/elevator.h
type Elevator struct {
	ID             int
	Floor          int
	MotorDirection MotorDirection
	MotorError	   bool
	Orders         [NumFloors][NumButtons]int
	Behaviour      ElevatorBehaviour
	Config         struct{
		ClearOrdersVariant ClearOrdersVariant
		DoorOpenDurationSec float64
	}
}

type NetOrder struct{
	ReceivingElevatorID int
	SendingElevatorID   int
	Order               ButtonEvent
	ReassignNum 		int
}

// Types for input to hall_request_assigner in distributor module
type ElevatorTagged struct  {
	Behaviour string                 `json:"behaviour"`
	Floor int                        `json:"floor"`
	MotorDirection string            `json:"direction"`
	CabOrders [NumFloors]bool `json:"cabRequests"`
}
type ElevatorsTagged struct{
	HallOrders [NumFloors][2]bool `json:"hallRequests"`
	States map[string]ElevatorTagged     `json:"states"`
}