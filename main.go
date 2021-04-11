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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	hardwareIO.Init(system.LocalHost, system.NumFloors)

	orderToSelfCh := make(chan system.ButtonEvent)
	hallOrderCh := make(chan system.ButtonEvent)
	floorArrivalCh := make(chan int)
	obstructionEventCh := make(chan bool)
	doorTimerDurationCh := make(chan float64)
	doorTimerTimedOutCh := make(chan bool)

	ownElevatorCh := make(chan system.Elevator)

	//Network channels for transmitting and receiving
	receiveElevatorInfoCh := make(chan system.ElevatorInformation)
	broadcastElevatorInfoCh := make(chan system.ElevatorInformation)
	receiveOrderCh := make(chan system.SendingOrder)
	sendPlacedMessageCh := make(chan system.SendingOrder)


	sendingOrderThroughNetCh := make(chan system.SendingOrder) //channel that receives
	elevatorInfoCh := make(chan system.ElevatorInformation) //channel with elevatorinformation, sent from networkmodule to
	//own modules

	messageTimerCh := make(chan system.SendingOrder)
	messageTimerTimedOutCh := make(chan system.SendingOrder)
	orderTimerCh := make(chan system.SendingOrder)
	orderTimerTimedOutCh := make(chan system.SendingOrder)

	// Timers:
	go timer.RunDoorTimer(doorTimerDurationCh, doorTimerTimedOutCh)
	go timer.RunMessageTimer(messageTimerCh, messageTimerTimedOutCh)
	go timer.RunOrderTimer(orderTimerCh, orderTimerTimedOutCh)

	// Hardware:
	go hardwareIO.RunHardware(orderToSelfCh, hallOrderCh, floorArrivalCh, obstructionEventCh)

	// Distributor and FSM:
	go fsm.ElevatorFSM(orderToSelfCh, floorArrivalCh, obstructionEventCh, ownElevatorCh, doorTimerDurationCh, doorTimerTimedOutCh)
	go distributor.OrderDistributor(hallOrderCh, elevatorInfoCh, ownElevatorCh, sendingOrderThroughNetCh, orderToSelfCh, messageTimerCh, messageTimerTimedOutCh, orderTimerCh, orderTimerTimedOutCh)

	//Network:
	go sendandreceive.GetReceiverAndTransmitterPorts(receiveElevatorInfoCh, broadcastElevatorInfoCh, receiveOrderCh, sendPlacedMessageCh, sendingOrderThroughNetCh, messageTimerCh)
	go sendandreceive.SendReceiveOrders(ownElevatorCh, broadcastElevatorInfoCh, receiveElevatorInfoCh, elevatorInfoCh, receiveOrderCh, orderToSelfCh, sendPlacedMessageCh)

	for {}
}