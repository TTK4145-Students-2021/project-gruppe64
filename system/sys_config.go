package system

// Project relevant configurations
const (
	NumElevators = 3
	ElevatorID = 0
	LocalHost = "localhost:15657"
)

// Other configurations
const (
	//Elevator
	NumFloors = 4
	NumButtons = 3
	ElevatorClearOrdersVariant = COInMotorDirection
	ElevatorDoorOpenDuration = 3.0 //Sec

	//HardwareIO
	CheckMotorAfterDuration = 3.0 //Sec

	//Distributor
	MaxReassignNum = 3

	//Timer
	MessageTimerDuration = 3 //Sec
	OrderTimerDuration = 20 //Sec

	//Network
	NetResendNum = 10
)