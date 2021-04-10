package sendandreceive

import (
	"fmt"
	"os"
	"realtimeProject/project-gruppe64/network/bcast"
	"realtimeProject/project-gruppe64/network/localip"
	"realtimeProject/project-gruppe64/system"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.


/*
func BroadcastElevator( e chan <- ElevatorInformation, elevator ElevatorInformation, buttonEvent chan <- hardwareIO.ButtonEvent) {
	for {

		time.Sleep(1 * time.Second)
		e <- elevator
	}
}
 */

//broadcastElevator sends the struct of the elevator to a channel.
func broadcastElevator(elevatorInfo chan <- system.ElevatorInformation,  elevator  system.Elevator) {
	//fsmElevator := <- elevator
	e := fsmElevatorToElevatorNetwork(elevator)
	elevatorInfo <- e
}

func sendOtherElevators(elevatorInfo chan <- system.ElevatorInformation,  otherElevator system.ElevatorInformation) {
	elevatorInfo <- otherElevator
}

func fsmElevatorToElevatorNetwork(elevator system.Elevator) system.ElevatorInformation {
	elevatorToSend := system.ElevatorInformation{}
	elevatorToSend.ID = system.ElevatorID
	elevatorToSend.Floor = elevator.Floor
	elevatorToSend.MotorDirection = elevator.MotorDirection
	elevatorToSend.Orders = elevator.Orders
	elevatorToSend.Behaviour = elevator.Behaviour
	return elevatorToSend
}

//assuming the order is consisting of int, int and [][]int where the first int is the ID of the elevator receiving
//the order, the second int is the ID of the elevator sending the order and the third is the actual order.
//SendOrder sends an order to another elevator 10 times.

//får sendingOrderThroughNet fra distributor, plasserer den i en annen nettverksmodul.
func sendOrder(placeOrder chan <- system.SendingOrder, sendingOrderThroughNet system.SendingOrder){
	for i := 0; i < 10; i++{
		time.Sleep(1 * time.Millisecond)
		placeOrder <- sendingOrderThroughNet
	}
}

/*
func sendAcceptMessage (acceptOrder chan <- SendingOrder, sendAcceptThroughNet SendingOrder){
	//orderToBeSent := <- a
	acceptOrder <- sendAcceptThroughNet //dette er som funksjonen sendAccept, så tror ikke vi trenger denne.
}
*/

func sendOrderMessage (orderToMessageTimer chan <- system.SendingOrder, acceptOrder chan <- system.SendingOrder, order system.SendingOrder){
	if order.SendingElevatorID == system.ElevatorID { //Fordi den sender vel ikke seg selv melding? Kanskje
		orderToMessageTimer <- order
		return
	}
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


func GetReceiverAndTransmitterPorts(othersElevatorInfo  chan system.ElevatorInformation, placedOrder chan system.SendingOrder,
	placeOrder chan system.SendingOrder, elevatorInfo chan system.ElevatorInformation){
	//peerGet chan peers.PeerUpdate, peerTXEnable chan bool, orderBack chan SendingOrder, orderBackSent chan SendingOrder){ //the two last ones are to check that SendAccept work

	//peerUpdate := UpdatePeer(fsm.ElevatorID)

	//Don't really think we need these two since we have IDs and ports but whatever
	//go peers.Receiver(50000, peerGet)
	//go peers.Transmitter(50000, peerUpdate, peerTXEnable)

	//these are to send and receive ElevatorStructs
	go bcast.Receiver(60000, othersElevatorInfo)
	go bcast.Transmitter(60000, elevatorInfo)

	//this is where the elevator will receive orders
	go bcast.Receiver(60001+system.ElevatorID, placeOrder)
	//go bcast.Receiver(60001+elevatorID+1, orderBackSent) //To test acceptMessage

	for elevatorIDs := 0; elevatorIDs < system.NumElevators; elevatorIDs ++{ //distributor.NumElevators+1; elevatorIDs++{ //Gjerne ha numElevators her.
		if elevatorIDs != system.ElevatorID { //bytt til == for å sjekke på samme node
			go bcast.Transmitter(60001 +system.ElevatorID, placedOrder)
			//go bcast.Transmitter(60001 + elevatorIDs+1, orderBack) //To test AcceptMessage

		}
	}

}
func SendReceiveOrders(elevator <- chan system.Elevator, otherElevatorInfo <- chan system.ElevatorInformation, //peerUpdate <- chan peers.PeerUpdate,
	sendingOrderThroughNet <- chan system.SendingOrder, placedOrder <- chan system.SendingOrder,
	elevatorInfo chan <- system.ElevatorInformation, placeOrder chan <- system.SendingOrder,
	acceptOrder chan <- system.SendingOrder, messageTimer chan <- system.SendingOrder) {
	for {
		select {
		/*
		case p := <-peerUpdate:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			*/

		case e := <-elevator:
			fmt.Printf("Elevatorinfo broadcasted: %#v\n", fsmElevatorToElevatorNetwork(e))
			broadcastElevator(elevatorInfo, e)

			case o := <-otherElevatorInfo:
				fmt.Printf("Elevatorinfo from other elevators and own broadcasted: %#v\n", o)
				sendOtherElevators(elevatorInfo, o)


		case s := <-sendingOrderThroughNet:
			fmt.Printf("Order sent through network: %#v\n", s)
			sendOrder(placeOrder, s)

		case p := <-placedOrder:
			fmt.Printf("Order info broadcasted:%#v\n", p)
			sendOrderMessage(messageTimer, acceptOrder, p)
		}
	}
}