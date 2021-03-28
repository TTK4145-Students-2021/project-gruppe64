package main

import (
	"./network/peers"
	"./network/sendandreceive"
	"./network/bcast"
	"fmt"
)

//ports
//19001: for IP-addresses to elevator1
////19002: for IP-addresses to elevator2
//19003: for IP-addresses to elevator3
//19998: for elevatorStructs to elevator1
//19999: for elevatorStructs to elevator2
//20000: for elevatorStructs to elevator3
//20001: for sending orders to elevator1
//20002: for sending orders to elevator2
//20003: for sending orders to elevator3

func main() {

	//test order
	testSlice := []interface{}{[]int{1}, []int{1}, []int {1,2}} //bør kanskje og ha hvilken knapp det er tryket på
	testSlice[0] = 1
	testSlice[1] = 2
	testSlice[2] = []int{1,2}

	testSlice2 := []interface{}{[]int{1}, []int{1}, []int {1,2}} //bør kanskje og ha hvilken knapp det er tryket på
	testSlice2[0] = 2
	testSlice2[1] = 1
	testSlice2[2] = []int{1,2}


	//test elevator
	elevator := sendandreceive.Elevator{ID: 1, MotorDirection: 3, Behaviour: "EB_idle"}
	elevator.Orders[1][1] = 1
	elevator.Orders[0][2] = 1


	peerUpdate := sendandreceive.UpdatePeer(elevator.ID)
	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(19001, peerUpdate, peerTxEnable)
	go peers.Receiver(19000, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	elevStructToSend := make(chan sendandreceive.Elevator)
	elevStructSent := make(chan sendandreceive.Elevator)

	orderToSend := make(chan []interface{})
	orderSent := make(chan []interface{})
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	//go sendandreceive.Transmit(elevator.ID, elevStructSent, orderSent, peerUpdateCh)

	//go bcast.Transmitter(19999, elevStructToSend)
	go bcast.Receiver(19998, elevStructSent) //og denne

	//go bcast.Transmitter(20002, orderToSend)
	go bcast.Receiver(20001, orderSent)

	go sendandreceive.BroadcastElevator(19999, elevStructToSend, elevator) //will broadcast elevatorstruct each second
	go sendandreceive.SendOrder(20002, orderToSend, testSlice)
	go sendandreceive.SendAccept(20002, orderToSend, testSlice2)
	fmt.Println("Started")
	sendandreceive.SendReceiveOrders(elevStructSent, orderSent, peerUpdateCh)

}

