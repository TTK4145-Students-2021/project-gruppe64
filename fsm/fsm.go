package fsm

import (
	"realtimeProject/project-gruppe64/hardwareIO"
)

func InitializeElevator(floorArrival <-chan int, elevatorUpdate chan<- Elevator){

	//elevatorUpdateCh := make(chan Elevator)
	//currentElevatorNOCh := make(chan Elevator)
	//currentElevatorFACh := make(chan Elevator)
	//currentElevatorOECh := make(chan Elevator)

	//runTimerCh := make(chan float64)
	//timedOutCh := make(chan bool)

	//go UpdateElevatorInformation(currentElevatorNOCh, currentElevatorFACh, currentElevatorOECh, elevatorUpdateCh)
	//go NewOrderFSM(buttonPressed, currentElevatorNOCh, elevatorUpdateCh, runTimerCh)
	//go FloorArrivalFSM(floorArrival, currentElevatorFACh, elevatorUpdateCh, runTimerCh)
	//go OrderExecutedFSM(timedOutCh, currentElevatorOECh, elevatorUpdateCh)

	//go timer.RunTimer(runTimerCh, timedOutCh)

	initializedElevator := Elevator{}
	select {
	case f :=<- floorArrival: // If the floor sensor registers a floor at initialization
		initializedElevator.Floor = f
		initializedElevator.MotorDirection = hardwareIO.MD_Stop
		initializedElevator.Behaviour = EB_Idle
		initializedElevator.Config.ClearOrdersVariant = CO_InMotorDirection
		initializedElevator.Config.DoorOpenDurationSec = 5.0
		elevatorUpdate <- initializedElevator
		break
	default: // If no floor is detected by the floor sensor
		initializedElevator.Floor = -1
		initializedElevator.MotorDirection = hardwareIO.MD_Down
		hardwareIO.SetMotorDirection(hardwareIO.MD_Down)
		initializedElevator.Behaviour = EB_Moving
		initializedElevator.Config.ClearOrdersVariant = CO_InMotorDirection
		initializedElevator.Config.DoorOpenDurationSec = 5.0
		elevatorUpdate <- initializedElevator
		break
	}
}


func UpdateElevatorInformation(currentElevatorNO chan<- Elevator, currentElevatorFA chan<- Elevator, currentElevatorOE chan<- Elevator, elevatorUpdate <-chan Elevator){
	for{
		select {
		case e := <- elevatorUpdate:
			printElevator(e)
			currentElevatorNO <- e
			currentElevatorFA <- e
			currentElevatorOE <- e
			break
		default:
			break
		}

	}
}

func setAllButtonLights(e Elevator){
	for f := 0; f < hardwareIO.NumFloors; f++ {
		for b := 0; b < hardwareIO.NumButtons; b++  {
			if e.Orders[f][b] != 0 {
				hardwareIO.SetButtonLamp(hardwareIO.ButtonType(b), f, true)
			} else {
				hardwareIO.SetButtonLamp(hardwareIO.ButtonType(b), f, false)
			}
		}
	}
}

func NewOrderFSM(btnEvent <-chan hardwareIO.ButtonEvent, currentElevator <-chan Elevator, elevatorUpdate chan<- Elevator, timerDur chan<- float64){
	localElevator := Elevator{}
	for{
		select {
		case e := <- currentElevator:
			localElevator = e
		case bE := <- btnEvent:
			hardwareIO.SetButtonLamp(bE.Button, bE.Floor, true)
			switch localElevator.Behaviour{
			case EB_DoorOpen:
				if localElevator.Floor == bE.Floor {
					timerDur <- localElevator.Config.DoorOpenDurationSec
				} else {
					localElevator.Orders[bE.Floor][int(bE.Button)] = 1
				}
				elevatorUpdate <- localElevator
				break
			case EB_Moving:
				localElevator.Orders[bE.Floor][int(bE.Button)] = 1
				elevatorUpdate <- localElevator
				break
			case EB_Idle:
				if localElevator.Floor == bE.Floor {
					hardwareIO.SetDoorOpenLamp(true)
					timerDur <- localElevator.Config.DoorOpenDurationSec
					localElevator.Behaviour = EB_DoorOpen
					elevatorUpdate <- localElevator
				} else {
					localElevator.Orders[bE.Floor][int(bE.Button)] = 1
					localElevator.MotorDirection = chooseDirection(localElevator)
					hardwareIO.SetMotorDirection(localElevator.MotorDirection)
					localElevator.Behaviour = EB_Moving
					elevatorUpdate <- localElevator
				}
				break
			}

		default:
			break
		}
	}
}


//NO
func FloorArrivalFSM(floor <-chan int, currentElevator <-chan Elevator, elevatorUpdate chan<- Elevator, timerDur chan<- float64){
	localElevator := Elevator{}
	for{
		select {
		case e := <-currentElevator:
			localElevator = e
		case f := <- floor:
			localElevator.Floor = f
			hardwareIO.SetFloorIndicator(localElevator.Floor)
			switch localElevator.Behaviour {
			case EB_Moving:
				if elevatorShouldStop(localElevator){
					hardwareIO.SetMotorDirection(hardwareIO.MD_Stop)
					hardwareIO.SetDoorOpenLamp(true)
					localElevator = clearOrdersAtCurrentFloor(localElevator)
					timerDur <- localElevator.Config.DoorOpenDurationSec
					setAllButtonLights(localElevator)
					localElevator.Behaviour = EB_DoorOpen
					elevatorUpdate <- localElevator
				} else if localElevator.Floor == 0{
					localElevator.MotorDirection = hardwareIO.MD_Up
					elevatorUpdate <- localElevator
				} else if localElevator.Floor == 3 {
					localElevator.MotorDirection = hardwareIO.MD_Down
					elevatorUpdate <- localElevator
				}

				break
			default:
				break
			}
			setAllButtonLights(localElevator)
		default:
			break
		}
	}
}

//OE
func OrderExecutedFSM(timedOut <- chan bool, currentElevator <-chan Elevator, elevatorUpdate chan<- Elevator){
	localElevator := Elevator{}
	for {
		select {
		case e := <-currentElevator:
			localElevator = e
		case t := <-timedOut:
			if t {
				switch localElevator.Behaviour {
				case EB_DoorOpen:
					clearOrdersAtCurrentFloor(localElevator)
					localElevator.MotorDirection = chooseDirection(localElevator)
					hardwareIO.SetDoorOpenLamp(false)
					hardwareIO.SetMotorDirection(localElevator.MotorDirection)
					if localElevator.MotorDirection == hardwareIO.MD_Stop {
						localElevator.Behaviour = EB_Idle
					} else {
						localElevator.Behaviour = EB_Moving
					}
					elevatorUpdate <- localElevator
					break
				default:
					break
				}
			}
		default:
			break
		}
	}
}