package main

import (
	"realtimeProject/project-gruppe64/distributor"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/network/sendandreceive"
	"realtimeProject/project-gruppe64/system"
	"realtimeProject/project-gruppe64/timer"
	"runtime"
)

/*
import (
	"./distributor"
	"./fsm"
	"./hardwareIO"
	"./network/peers"
	"./network/sendandreceive"
	"./system"
	"./timer"
	"io/ioutil"
	"runtime"
	"strconv"
	"time"
)
*/

func primaryWork(activateAsPrimary <-chan bool){
	for{
		select{
		case activate := <-activateAsPrimary:
			if activate{
				hardwareIO.Init(system.LocalHost, system.NumFloors)
				system.SpawnBackup()
				go system.PrimaryDocumentation()

				// ->FSM
				floorArrivalCh := make(chan int)
				obstructionEventCh := make(chan bool)
				orderToSelfCh := make(chan system.ButtonEvent)
				doorTimerTimedOutCh := make(chan bool)

				// ->Distributor
				hallOrderCh := make(chan system.ButtonEvent)
				ownElevatorCh := make(chan system.Elevator)
				otherElevatorCh := make(chan system.Elevator)
				messageTimerTimedOutCh := make(chan system.NetOrder)
				orderTimerTimedOutCh := make(chan system.NetOrder)
				elevatorConnectedCh := make(chan int, system.NumElevators - 1)
				elevatorDisconnectedCh := make(chan int, system.NumElevators - 1)

				// ->Network
				shareOwnElevatorCh := make(chan system.Elevator)
				orderThroughNetCh := make(chan system.NetOrder)


				// ->Timer
				doorTimerDurationCh := make(chan float64)
				messageTimerCh := make(chan system.NetOrder)
				placedMessageReceivedCh := make(chan system.NetOrder)
				orderTimerCh := make(chan system.NetOrder)

				// Hardware:
				go hardwareIO.RunHardware(orderToSelfCh, hallOrderCh, floorArrivalCh, obstructionEventCh)

				// FSM:
				go distributor.OrderDistributor(hallOrderCh, otherElevatorCh, ownElevatorCh, shareOwnElevatorCh,
					orderThroughNetCh, orderToSelfCh, messageTimerCh, messageTimerTimedOutCh, orderTimerCh,
					orderTimerTimedOutCh, elevatorConnectedCh, elevatorDisconnectedCh)

				// Distributor:
				go fsm.ElevatorFSM(orderToSelfCh, floorArrivalCh, obstructionEventCh, ownElevatorCh,
					doorTimerDurationCh, doorTimerTimedOutCh)

				// Network:
				go sendandreceive.RunNetworking(shareOwnElevatorCh, otherElevatorCh, orderThroughNetCh,
					placedMessageReceivedCh, orderToSelfCh, elevatorConnectedCh, elevatorDisconnectedCh)

				// Timers:
				go timer.RunDoorTimer(doorTimerDurationCh, doorTimerTimedOutCh)
				go timer.RunMessageTimer(messageTimerCh, placedMessageReceivedCh, messageTimerTimedOutCh)
				go timer.RunOrderTimer(orderTimerCh, orderTimerTimedOutCh)

			}
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	activateAsPrimaryCh := make(chan bool)
	if system.IsBackup(){
		go system.CheckPrimaryExistence(activateAsPrimaryCh)
		go primaryWork(activateAsPrimaryCh)
	} else {
		system.MakeBackupFile()
		go primaryWork(activateAsPrimaryCh)
		activateAsPrimaryCh <- true
	}
	for {}
}