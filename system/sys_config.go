package system

const (
	NumElevators = 3
	ElevatorID =1

	LocalHost = "localhost:15661"
)

//Other configurations
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