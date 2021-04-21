package system

type MotorDirection int
const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int
const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

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
}

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