package messages

import (
	"../../../../fsm2"
	"../../../../hardwareIO"
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
	Order hardwareIO.ButtonEvent
}


//Denne funker, men må inkludere hardwareIO
type ElevatorInformation struct {
	ID             int
	Floor          int
	MotorDirection hardwareIO.MotorDirection
	Orders         [hardwareIO.NumFloors][hardwareIO.NumButtons]int
	Behaviour      fsm2.ElevatorBehaviour
}

//BroadcastElevator sends the struct of the elevator to a channel.
func BroadcastElevator( e chan <- ElevatorInformation, elevator ElevatorInformation) {
	for {
		time.Sleep(1 * time.Second)
		e <- elevator
	}
}
func FSMElevatorToElevatorNetwork(elevator fsm2.Elevator, ID int) ElevatorInformation{
	elevatorToSend := ElevatorInformation{}
	elevatorToSend.ID = ID
	elevatorToSend.Floor = elevator.Floor
	elevatorToSend.MotorDirection = elevator.MotorDirection
	elevatorToSend.Orders = elevator.Orders
	elevatorToSend.Behaviour = elevator.Behaviour
	return elevatorToSend
}
//assuming the order is consisting of int, int and [][]int where the first int is the ID of the elevator receiving
//the order, the second int is the ID of the elevator sending the order and the third is the actual order.
//SendOrder sends an order to another elevator 10 times.
func SendOrder(placeOrder chan <- OrderToSend, order OrderToSend, update <- chan bool){
	for i := 0; i < 10; i++{
		time.Sleep(1 * time.Second)
		placeOrder <- order
		//turn into milliseconds or less when actually using it
	}
}

//SendAccept sends the order back to the elevator, to let it know that it takes the order.
func SendAccept(acceptOrder chan <- OrderToSend, order OrderToSend) {
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
/*
func GetReceiverAndTransmitterPorts(elevatorID int, elevStructSent chan Elevator, orderSent chan OrderToSend, peerGet chan peers.PeerUpdate,
	orderToSend chan OrderToSend, elevatorToSend chan Elevator, peerTXEnable chan bool){
	peerUpdate := UpdatePeer(elevatorID)
	go peers.Receiver(19000, peerGet)
	go peers.Transmitter(19000, peerUpdate, peerTXEnable)
	go bcast.Receiver(20000, elevStructSent)
	go bcast.Transmitter(20000, elevatorToSend)
	go bcast.Transmitter(1 + 20000, orderToSend) //

	//go bcast.Receiver(elevatorID + 19997, elevStructSent)
	go bcast.Receiver(elevatorID + 20000, orderSent)
}
*/


func GetReceiverAndTransmitterPorts(elevatorID int, elevStructSent chan ElevatorInformation, orderSent chan OrderToSend, peerGet chan peers.PeerUpdate,
	orderToSend chan OrderToSend, elevatorToSend chan ElevatorInformation, peerTXEnable chan bool,
	orderBack chan OrderToSend, orderBackSent chan OrderToSend){ //the two last ones are to check that SendAccept work
	peerUpdate := UpdatePeer(elevatorID)

	//Don't really think we need these two since we have IDs and ports but whatever
	go peers.Receiver(19000, peerGet)
	go peers.Transmitter(19000, peerUpdate, peerTXEnable)

	//these are to send and receive ElevatorStructs
	go bcast.Receiver(20000, elevStructSent)
	go bcast.Transmitter(20000, elevatorToSend)

	//this is where the elevator will receive orders
	go bcast.Receiver(20000+elevatorID, orderSent)
	go bcast.Receiver(20000+elevatorID+1, orderBackSent) //To test acceptMessage




	for elevatorIDs := 1; elevatorIDs < 4; elevatorIDs++{ //Gjerne ha numElevators her.
		if elevatorIDs == elevatorID{ //bytt til == for å sjekke på samme node
			go bcast.Transmitter(20000 + elevatorIDs, orderToSend)
			go bcast.Transmitter(20000 + elevatorIDs+1, orderBack) //To test AcceptMessage

		}
	}

}

//SendReceiveOrders states what happens when orders are sent and so on. This function is to be changed when everything
//is finished of course, so that the different modules communicate with each other.
func SendReceiveOrders(elevStructSent chan  ElevatorInformation, orderSent chan OrderToSend,
	peerUpdate <- chan peers. PeerUpdate, orderBack chan OrderToSend, orderBackSent chan OrderToSend){
	//to siste er for å sjekke AcceptMessage
	for {
		select {
		case p := <-peerUpdate:
			//should save the different IP-adresses and elevatorIDs in some way. IP is not really needed since we have
			//designated listeningports
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
			SendAccept(orderBack, b)


		case c := <-orderBackSent:
			fmt.Printf("The order: %#v was sent back!\n", c)

		}
	}
}
