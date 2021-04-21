package sendandreceive

import (
	"fmt"
	"realtimeProject/project-gruppe64/network/bcast"
	"realtimeProject/project-gruppe64/network/peers"
	"realtimeProject/project-gruppe64/system"
	"strconv"
	"time"
)

/*
import (
	"../../system"
	"../bcast"
	"../peers"
	"fmt"
	"strconv"
	"time"
)
*/

func RunNetworking(shareOwnElevatorCh <-chan system.Elevator, otherElevatorCh chan<- system.Elevator,
	orderThroughNetCh <-chan system.NetOrder, placedMessageReceievedCh chan<- system.NetOrder,
	orderToSelfCh chan<- system.ButtonEvent, elevatorConnectedCh chan<- int, elevatorDisconnectedCh chan<- int) {

	transmitPeerBoolCh := make(chan bool)
	receivePeerCh := make(chan peers.PeerUpdate)

	receiveElevatorCh := make(chan system.Elevator)
	transmitElevatorCh := make(chan system.Elevator)

	receiveOrderCh := make(chan system.NetOrder)  //Used for both placed message and orders


	go peers.Transmitter(59999, strconv.Itoa(system.ElevatorID), transmitPeerBoolCh)
	go peers.Receiver(59999, receivePeerCh)

	go bcast.Receiver(60000, receiveElevatorCh)
	go bcast.Transmitter(60000, transmitElevatorCh)

	go bcast.Receiver(60001+system.ElevatorID, receiveOrderCh)

	go elevatorsShareNet(shareOwnElevatorCh, transmitElevatorCh, receiveElevatorCh, otherElevatorCh)
	go peersNet(receivePeerCh, elevatorConnectedCh, elevatorDisconnectedCh)

	for elevID := 0; elevID < system.NumElevators; elevID++ {
		if elevID != system.ElevatorID {
			transmitOrderCh := make(chan system.NetOrder) //Reset every run
			go bcast.Transmitter(60001 +elevID, transmitOrderCh) //Transmit orders to place
			go ordersNet(elevID, orderThroughNetCh, placedMessageReceievedCh, transmitOrderCh, receiveOrderCh, orderToSelfCh)
		}
	}
}

func elevatorsShareNet(shareOwnElevatorCh <- chan system.Elevator, transmitElevatorCh chan <- system.Elevator, receiveElevatorCh <- chan system.Elevator,
	otherElevatorCh chan<- system.Elevator) {
	for {
		select {
		case shareOwnElevator := <-shareOwnElevatorCh:
			//fmt.Printf("Elevatorinfo broadcasted: %#v\n", system.ElevatorInformation{ID: system.ElevatorID, Floor: ownElevator.Floor, MotorDirection: ownElevator.MotorDirection, Orders: ownElevator.Orders, Behaviour: ownElevator.Behaviour})
			transmitElevatorCh <- shareOwnElevator

		case receiveElevator := <-receiveElevatorCh:
			if receiveElevator.ID != system.ElevatorID {
				//fmt.Printf("Elevatorinfo from other elevator: %#v\n", receiveElevator)
				otherElevatorCh <- receiveElevator
			}
		}
	}
}

func peersNet(receivePeerCh <-chan peers.PeerUpdate, elevatorConnectedCh chan<- int, elevatorDisconnectedCh chan<- int){
	for {
		select {
		case receivePeer := <-receivePeerCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", receivePeer.Peers)
			fmt.Printf("  New:      %q\n", receivePeer.New)
			fmt.Printf("  Lost:     %q\n", receivePeer.Lost)

			if receivePeer.New != "" && receivePeer.New != strconv.Itoa(system.ElevatorID) {
				fmt.Println("New peer ID: " + receivePeer.New)
				newSentID, _ := strconv.Atoi(receivePeer.New)
				elevatorConnectedCh <- newSentID
			}
			for IDLost := 0; IDLost < len(receivePeer.Lost); IDLost ++{
				lostSentID,_ := strconv.Atoi(receivePeer.Lost[IDLost])
				fmt.Println("Lost sent ID:", lostSentID)
				elevatorDisconnectedCh <- lostSentID
			}
		}
	}
}

func ordersNet(threadElevatorID int, sendingOrderThroughNet <-chan system.NetOrder, placedMessageRecieved chan<- system.NetOrder, networkSend chan<- system.NetOrder, networkReceive <-chan system.NetOrder, orderToSelf chan<- system.ButtonEvent) {
	for {
		select {
		case sOrdNet := <-sendingOrderThroughNet:
			if sOrdNet.ReceivingElevatorID == threadElevatorID {
				fmt.Printf("Order sent through network: %#v\n", sOrdNet)
				for i := 0; i < system.NetResendNum; i++ {
					time.Sleep(1 * time.Millisecond)
					networkSend <- sOrdNet
				}
			}
		case netReceive := <-networkReceive:
			if netReceive.SendingElevatorID == system.ElevatorID { //THEN IT IS A PLACED MESSAGE
				fmt.Println("Placed message reveived")
				placedMessageRecieved <- netReceive
			}

			if netReceive.ReceivingElevatorID == system.ElevatorID { //THEN IT IS A ORDER
				fmt.Printf("Order received: %#v\n", netReceive)
				orderToSelf <- netReceive.Order
				for i := 0; i < system.NetResendNum; i++ {
					networkSend <- netReceive //As placed message }
				}
			}
		}
	}
}
