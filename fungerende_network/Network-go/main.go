
package main

import (
	//"Network-go/network/bcast"
	"./network/peers"
	"./network/sendandreceive"
	"fmt"
)

//ports
//19000: for IP-addresses
//20000: for elevatorStructs
//20001: for sending orders to elevator1
//20002: for sending orders to elevator2
//20003: for sending orders to elevator3

func main() {
	//test order
	testOrder := sendandreceive.OrderToSend{ReceivingElevatorID: 2, SendingElevatorID: 1}
	testOrder.Order[0] = 1
	testOrder.Order[1] = 3

	//test elevator
	elevator := sendandreceive.Elevator{ID: 2, MotorDirection: 1, Behaviour: "EB_idle"}
	elevator.Orders[2][3] = 1
	elevator.Orders[1][1] = 1

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	// We make channels for sending and receiving our custom data types
	elevStructToSend := make(chan sendandreceive.Elevator)
	elevStructSent := make(chan sendandreceive.Elevator)

	orderToSend := make(chan sendandreceive.OrderToSend)
	orderSent := make(chan sendandreceive.OrderToSend)

	go sendandreceive.GetReceiverAndTransmitterPorts(elevator.ID, elevStructSent, orderSent, peerUpdateCh, orderToSend, elevStructToSend, peerTxEnable)
	go sendandreceive.BroadcastElevator(elevStructToSend, elevator) //will broadcast elevatorstruct each second
	go sendandreceive.SendOrder(orderToSend, testOrder, peerTxEnable)
	fmt.Println("Started")
	sendandreceive.SendReceiveOrders(elevStructSent, orderSent, orderSent, peerUpdateCh)
}
