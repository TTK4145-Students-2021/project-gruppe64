package sendandreceive

import (
	"../localip"
	"../peers"
	"../bcast"
	"fmt"
	"os"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.

type OrderToSend struct{
	ReceivingElevatorID int
	SendingElevatorID int
	Order [2] int
}

type Elevator struct {
	ID     int
	Floor int
	MotorDirection int 		//egentlig hardwareIO.MotorDirection
	Orders [3][4] int 		// egentlig [hardwareIO.NumFloors][hardwareIO.NumButtons] int
	Behaviour string 		 //egentlig ElevatorBehaviour
	//tenker og at vi kanskje kan drite i å sende configStructen
}

//BroadcastElevator sends the struct of the elevator to a channel.
func BroadcastElevator( e chan <- Elevator, elevator Elevator) {
	
	for {
		time.Sleep(1 * time.Second)
		e <- elevator
	}
}

//assuming the order is consisting of int, int and [][]int where the first int is the ID of the elevator receiving
//the order, the second int is the ID of the elevator sending the order and the third is the actual order.
//SendOrder sends an order to another elevator 10 times.
func SendOrder(placeOrder chan <- OrderToSend, order OrderToSend, update <- chan bool){
	/*elevatorReceiver := order[0].(int)
	port := 20000 + elevatorReceiver
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	bcast.Transmitter(port, localIP, update)
	*/

	for i := 0; i < 10; i++{
		time.Sleep(1 * time.Second)
		placeOrder <- order
		 //turn into milliseconds or less when actually using it
	}
}

//SendAccept sends the order back to the elevator, to let it know that it takes the order.
func SendAccept(acceptOrder chan <- OrderToSend, order OrderToSend) {
	//elevatorReceiver := order[1].(int)
	/*
	port := 20000 + elevatorReceiver
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	*/
	//go bcast.Transmitter(port, localIP, update)
	acceptOrder <- order
}


//UpdatePeer returns a string consisting of the elevatorID and the IP-address of the elevator node. This will then be
//sent to a dedicated port, 19000.
func UpdatePeer(elevatorID int) string{
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
		eIDeIP := fmt.Sprintf("Elevator ID: %d, IP: %s-%d", elevatorID, localIP, os.Getpid()) //not sure if we need os.Getpid() as well
		return eIDeIP
	}
	return fmt.Sprintf("Elevator ID: %d, IP: %s-%d", elevatorID, localIP, os.Getpid())
}

func GetReceiverAndTransmitterPorts(elevatorID int, elevStructSent chan Elevator, orderSent chan OrderToSend, peerGet chan peers.PeerUpdate, 
	orderToSend chan OrderToSend, elevatorToSend chan Elevator, peerTXEnable chan bool){
	peerUpdate := UpdatePeer(elevatorID)
	go peers.Receiver(19000, peerGet)
	go peers.Transmitter(19000, peerUpdate, peerTXEnable)
	go bcast.Receiver(20000, elevStructSent)
	go bcast.Transmitter(20000, elevatorToSend)
	go bcast.Transmitter(2 + 20000, orderToSend) //

	//go bcast.Receiver(elevatorID + 19997, elevStructSent)
	go bcast.Receiver(elevatorID + 20000, orderSent)
}

//SendReceiveOrders states what happens when orders are sent and so on. This function is to be changed when everything
//is finished of course, so that the different modules communicate with each other.
func SendReceiveOrders(elevStructSent chan  Elevator, orderToSend chan OrderToSend, orderSent chan OrderToSend, peerUpdate <- chan peers. PeerUpdate){
	for {
		select {
		case p := <-peerUpdate:
			//should save the different IP-adresses and elevatorIDs in some way. Maybe IP is not needed since we have
			//designated listeningports? Don't know.
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New) //don't feel like this is needed, so I removed it
			//fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-elevStructSent:
			//a should be saved in order designator module
			fmt.Printf("Received: %#v\n", a) //removing # returns just the values and not the types

		case b := <- orderSent:
			//b should be saved in order designator module
			fmt.Printf("Received: %#v\n", b)
			//SendAccept(orderSent, b)


		}
	}
}