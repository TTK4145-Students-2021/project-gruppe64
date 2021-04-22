package sendandreceive

import (
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/network/bcast"
	"realtimeProject/project-gruppe64/network/peers"
	"realtimeProject/project-gruppe64/system"
	"strconv"
	"time"
)


/*
import (
	"../../hardwareIO"
	"../../system"
	"../bcast"
	"../peers"
	"strconv"
	"time"
)
 */


func RunNetworking(shareOwnElevatorCh <-chan system.Elevator, otherElevatorCh chan<- system.Elevator,
	orderThroughNetCh <-chan system.NetOrder, placedMessageReceievedCh chan<- system.NetOrder,
	orderTimerCh chan<- system.NetOrder, orderToSelfCh chan<- system.ButtonEvent, elevatorConnectedCh chan<- int,
	elevatorDisconnectedCh chan<- int) {

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
			go ordersNet(elevID, orderThroughNetCh, placedMessageReceievedCh, orderTimerCh, transmitOrderCh,
				receiveOrderCh, orderToSelfCh)
		}
	}
}

func elevatorsShareNet(shareOwnElevatorCh <- chan system.Elevator, transmitElevatorCh chan <- system.Elevator, receiveElevatorCh <- chan system.Elevator,
	otherElevatorCh chan<- system.Elevator) {
	for {
		select {
		case shareOwnElevator := <-shareOwnElevatorCh:
			transmitElevatorCh <- shareOwnElevator

		case receiveElevator := <-receiveElevatorCh:
			if receiveElevator.ID != system.ElevatorID {
				otherElevatorCh <- receiveElevator
			}
		}
	}
}

func peersNet(receivePeerCh <-chan peers.PeerUpdate, elevatorConnectedCh chan<- int, elevatorDisconnectedCh chan<- int){
	for {
		select {
		case receivePeer := <-receivePeerCh:
			if receivePeer.New != "" && receivePeer.New != strconv.Itoa(system.ElevatorID) {
				newSentID, _ := strconv.Atoi(receivePeer.New)
				elevatorConnectedCh <- newSentID
			}
			for IDLost := 0; IDLost < len(receivePeer.Lost); IDLost ++{
				lostSentID,_ := strconv.Atoi(receivePeer.Lost[IDLost])
				elevatorDisconnectedCh <- lostSentID
			}
		}
	}
}

func ordersNet(threadElevatorID int, orderThroughNetCh <-chan system.NetOrder,
	placedMessageRecievedCh chan<- system.NetOrder, orderTimerCh chan<- system.NetOrder,
	transmitOrderCh chan<- system.NetOrder, receiveOrderCh <-chan system.NetOrder,
	orderToSelfCh chan<- system.ButtonEvent) {
	for {
		select {
		case orderThroughNet := <-orderThroughNetCh:
			if orderThroughNet.ReceivingElevatorID == threadElevatorID {
				for i := 0; i < system.NetResendNum; i++ {
					time.Sleep(1 * time.Millisecond)
					transmitOrderCh <- orderThroughNet
				}
			}
		case receiveOrder := <-receiveOrderCh:
			hardwareIO.SetButtonLamp(receiveOrder.Order.Button, receiveOrder.Order.Floor, true)
			if receiveOrder.SendingElevatorID == system.ElevatorID { // If true: is placed message
				placedMessageRecievedCh <- receiveOrder
				orderTimerCh <- receiveOrder
			}

			if receiveOrder.ReceivingElevatorID == system.ElevatorID { // If true: is order
				orderToSelfCh <- receiveOrder.Order
				for i := 0; i < system.NetResendNum; i++ {
					transmitOrderCh <- receiveOrder // As placed message
				}
			}
		}
	}
}