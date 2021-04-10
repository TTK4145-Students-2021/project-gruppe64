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
	hardwareIO.Init("localhost:15657", system.NumFloors)

	orderToSelfCh := make(chan system.ButtonEvent)
	hallOrderCh := make(chan system.ButtonEvent)
	floorArrivalCh := make(chan int)
	obstructionEventCh := make(chan bool)
	doorTimerDurationCh := make(chan float64)
	doorTimerTimedOutCh := make(chan bool)

	ownElevatorCh := make(chan system.Elevator)

	//Network channels
	sendingOrderThroughNetCh := make(chan system.SendingOrder) //channel that receives
	placedOrderCh := make(chan system.SendingOrder) //output from another network module into other network module
	elevatorInfoCh := make(chan system.ElevatorInformation) //channel with elevatorinformation, sent from networkmodule to
	//own modules
	othersElevatorInfoCh := make(chan system.ElevatorInformation)
	placeOrderCh := make(chan system.SendingOrder) //sent from this networkmodule to other network module
	acceptOrderCh := make(chan system.SendingOrder) //sent from this networkmodule to other network module

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
	go sendandreceive.GetReceiverAndTransmitterPorts(othersElevatorInfoCh, placedOrderCh, placeOrderCh, elevatorInfoCh)
	go sendandreceive.SendReceiveOrders(ownElevatorCh, othersElevatorInfoCh, sendingOrderThroughNetCh, placedOrderCh, elevatorInfoCh, placeOrderCh, acceptOrderCh, messageTimerCh)

	for {}
}