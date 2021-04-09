package main

import (
	"realtimeProject/distributor"
	"realtimeProject/fsm"
	"realtimeProject/hardwareIO"
	"realtimeProject/timer"
	"realtimeProject/Network-go/network/sendandreceive"
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


	//Network channels
	sendingOrderThroughNetCh := make(chan sendandreceive.SendingOrder) //channel that receives
	placedOrderCh := make(chan sendandreceive.SendingOrder) //output from another network module into other network module
	elevatorInfoCh := make(chan sendandreceive.ElevatorInformation) //channel with elevatorinformation, sent from networkmodule to
	//own modules
	othersElevatorInfoCh := make(chan sendandreceive.ElevatorInformation)
	placeOrderCh := make(chan sendandreceive.SendingOrder) //sent from this networkmodule to other network module
	acceptOrderCh := make(chan sendandreceive.SendingOrder) //sent from this networkmodule to other network module


	messageTimerCh := make(chan sendandreceive.SendingOrder)
	messageTimerTimedOutCh := make(chan sendandreceive.SendingOrder)
	orderTimerCh := make(chan sendandreceive.SendingOrder)
	orderTimerTimedOutCh := make(chan sendandreceive.SendingOrder)

	// Timers:
	go timer.RunDoorTimer(doorTimerDurationCh, doorTimerTimedOutCh)
	go timer.RunMessageTimer(messageTimerCh, messageTimerTimedOutCh)
	go timer.RunOrderTimer(orderTimerCh, orderTimerTimedOutCh)

	// Hardware:
	go hardwareIO.RunHardware(orderToSelfCh, hallOrderCh, floorArrivalCh, obstructionEventCh)
	go fsm.ElevatorFSM(orderToSelfCh, floorArrivalCh, obstructionEventCh, ownElevatorCh, doorTimerDurationCh, doorTimerTimedOutCh)
	go distributor.OrderDistributor(hallOrderCh, elevatorInfoCh, ownElevatorCh, sendingOrderThroughNetCh, orderToSelfCh, messageTimerCh, messageTimerTimedOutCh, orderTimerCh, orderTimerTimedOutCh)
	go sendandreceive.GetReceiverAndTransmitterPorts(othersElevatorInfoCh, placedOrderCh, placeOrderCh, elevatorInfoCh)
	go sendandreceive.SendReceiveOrders(ownElevatorCh, othersElevatorInfoCh, sendingOrderThroughNetCh, placedOrderCh, elevatorInfoCh, placeOrderCh, acceptOrderCh, messageTimerCh)
	for {}
}