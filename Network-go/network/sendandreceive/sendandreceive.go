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

type Elevator struct {
	ID     int
	Floor int
	MotorDirection int 		//egentlig hardwareIO.MotorDirection
	Orders [3][4] int 		// egentlig [hardwareIO.NumFloors][hardwareIO.NumButtons] int
	Behaviour string 		 //egentlig ElevatorBehaviour
	//tenker og at vi kanskje kan drite i å sende configStructen
}

//BroadcastElevator sends the struct of the elevator to a channel.
func BroadcastElevator(port int, elevatorChan chan Elevator, elevator Elevator) {
	go bcast.Transmitter(port, elevatorChan)
	for {
		time.Sleep(1 * time.Second)
		elevatorChan <- elevator
	}
}

//assuming the order is consisting of int, int and [][]int where the first int is the ID of the elevator receiving
//the order, the second int is the ID of the elevator sending the order and the third is the actual order.
//SendOrder sends an order to another elevator 10 times.
func SendOrder(port int, placeOrder chan []interface{}, order []interface{}){
	go bcast.Transmitter(port, placeOrder)
	for i := 0; i < 10; i++{
		placeOrder <- order
		time.Sleep(1 * time.Second) //turn into milliseconds or less when actually using it
	}
}

//SendAccept sends the order back to the elevator, to let it know that it takes the order.
func SendAccept(port int, acceptOrder chan []interface{}, order []interface{}) {
	go bcast.Transmitter(port, acceptOrder)
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

func Transmit(elevatorID int, elevChan chan  Elevator, orderChan chan [] interface{},peersChan chan peers.PeerUpdate){

}

//SendReceiveOrders states what happens when orders are sent and so on. This function is to be changed when everything
//is finished of course, so that the different modules communicate with each other.
func SendReceiveOrders( elevStructSent chan  Elevator, orderSent chan []interface{}, peerUpdate <- chan peers. PeerUpdate){
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

			//SendAccept(orderToSend, b)
			//fmt.Printf("Sent accept %#v\n", b)



		}
	}
}