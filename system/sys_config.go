package system

const (
	NumFloors = 4
	NumElevators = 3
	ElevatorID =2
	LocalHost = "localhost:15662"
)

//Other configurations
const (
	//Elevator
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