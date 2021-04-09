package main

import (
	"realtimeProject/project-gruppe64/distributor"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/timer"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	hardwareIO.Init("localhost:15657", hardwareIO.NumFloors)

	orderToSelfCh := make(chan hardwareIO.ButtonEvent)
	hallOrderCh := make(chan hardwareIO.ButtonEvent)
	floorArrivalCh := make(chan int)
	obstructionEventCh := make(chan bool)
	doorTimerDurationCh := make(chan float64)
	doorTimerTimedOutCh := make(chan bool)

	ownElevatorCh := make(chan fsm.Elevator)

	// skal være network... (typene definert der, ikke i distributor)
	elevatorInfoCh := make(chan distributor.ElevatorInformation)
	sendingOrderThroughNetCh := make(chan distributor.SendingOrder)

	messageTimerCh := make(chan distributor.SendingOrder)
	messageTimerTimedOutCh := make(chan distributor.SendingOrder)
	orderTimerCh := make(chan distributor.SendingOrder)
	orderTimerTimedOutCh := make(chan distributor.SendingOrder)

	// Timers:
	go timer.RunDoorTimer(doorTimerDurationCh, doorTimerTimedOutCh)
	go timer.RunMessageTimer(messageTimerCh, messageTimerTimedOutCh)
	go timer.RunOrderTimer(orderTimerCh, orderTimerTimedOutCh)

	// Hardware:
	go hardwareIO.RunHardware(orderToSelfCh, hallOrderCh, floorArrivalCh, obstructionEventCh)

	go fsm.ElevatorFSM(orderToSelfCh, floorArrivalCh, obstructionEventCh, ownElevatorCh, doorTimerDurationCh, doorTimerTimedOutCh)
	go distributor.OrderDistributor(hallOrderCh, elevatorInfoCh, ownElevatorCh, sendingOrderThroughNetCh, orderToSelfCh, messageTimerCh, messageTimerTimedOutCh, orderTimerCh, orderTimerTimedOutCh)

	for {}
}