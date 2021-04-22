package system

type MotorDirection int
const (
	MDUp   MotorDirection = 1
	MDDown                = -1
	MDStop                = 0
)

type ButtonType int
const (
	BTHallUp   ButtonType = 0
	BTHallDown            = 1
	BTCab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type ElevatorBehaviour int
const (
	EBIdle     ElevatorBehaviour = 0
	EBDoorOpen                   = 1
	EBMoving                     = 2
)

type ClearOrdersVariant int
const (
	COAll    ClearOrdersVariant = 0
	COInMotorDirection          = 1
)

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


// Types for for input to the distributor hall_request_assigner
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