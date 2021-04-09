
package main

import (
	"realtimeProject/Network-go/network/peers"
	"realtimeProject/Network-go/network/sendandreceive"
	//"realtimeProject/distributor"
	"realtimeProject/fsm"
	//"realtimeProject/hardwareIO"

	"flag"
	"fmt"
)

//ports
//19000: for IP-addresses
//20000: for elevatorStructs
//20001: for sending orders to elevator1
//20002: for sending orders to elevator2
//20003: for sending orders to elevator3

func main() {
	//setting elevator ID:
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	//var idInt int
	//flag.IntVar(&idInt, "idInt", 0, "id int of this peer")
	//flag.Parse()
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id` VERY IMPORTANT TO WRITE -id=our_id NOT -id = our_id
	// for idInt: -id int=our_id

	//test order
	//buttonEvent := hardwareIO.ButtonEvent{Floor:2,Button:1}
	//testOrder := sendandreceive.SendingOrder{ReceivingElevatorID: 2, SendingElevatorID: distributor.ElevatorID, Order: buttonEvent}

	//test elevator

	fsm2Elevator := fsm.Elevator{Floor: 2, MotorDirection: 2,Behaviour: fsm.EB_Moving}
	fsm2Elevator.Orders[2][0] = 1
	fsm2Elevator.Orders[0][0] = 1
	fsm2Elevator.Orders[1][0] = 1
	fsm2Elevator.Orders[3][0] = 1

	//elevator := sendandreceive.FSMElevatorToElevatorNetwork(fsm2Elevator, distributor.ElevatorID)
	//fmt.Printf("%#v", elevator)

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	elevStructToSend := make(chan sendandreceive.ElevatorInformation)
	elevStructSent := make(chan sendandreceive.ElevatorInformation)
	orderToSend := make(chan sendandreceive.SendingOrder)
	orderSent := make(chan sendandreceive.SendingOrder)

	orderBack := make(chan sendandreceive.SendingOrder)
	orderBackSent := make(chan sendandreceive.SendingOrder)
	//før var id elevator.ID
	go sendandreceive.GetReceiverAndTransmitterPorts(id, elevStructSent, orderSent, peerUpdateCh, orderToSend, elevStructToSend, peerTxEnable, orderBack, orderBackSent)
	fmt.Println("Started")
	//sendandreceive.SendReceiveOrders(elevStructSent, orderSent, peerUpdateCh, orderBack, orderBackSent)
}
