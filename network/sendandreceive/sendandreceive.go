package sendandreceive

import (
	"fmt"
	"realtimeProject/project-gruppe64/network/bcast"
	"realtimeProject/project-gruppe64/network/peers"
	"realtimeProject/project-gruppe64/system"
	"strconv"
	"time"
)

const (
	resendNum = 10
	)



func SetUpReceiverAndTransmitterPorts(receiveElevatorInfo chan system.ElevatorInformation, broadcastElevatorInfo chan system.ElevatorInformation,
	networkReceive chan system.SendingOrder, sendingOrderThroughNet <-chan system.SendingOrder, placedMessageReceived chan<- system.SendingOrder,
	orderToSelf chan<- system.ButtonEvent, receivePeers chan peers.PeerUpdate){

	transmitPeerBoolCh := make(chan bool)

	go peers.Transmitter(59999, strconv.Itoa(system.ElevatorID), transmitPeerBoolCh)
	go peers.Receiver(59999, receivePeers)

	go bcast.Receiver(60000, receiveElevatorInfo) //Receive others elevator information
	go bcast.Transmitter(60000, broadcastElevatorInfo) //Send elevator Information

	 //PLACED MESSAGE AND ORDER ON SAME. TO FROM/
	go bcast.Receiver(60001+system.ElevatorID, networkReceive) //Receive orders

	for elevID := 0; elevID < system.NumElevators; elevID++ {
		if elevID != system.ElevatorID {
			networkSendCh := make(chan system.SendingOrder) //Reset every run
			go bcast.Transmitter(60001 +elevID, networkSendCh) //Transmit orders to place
			go placeOrderNetworking(elevID, sendingOrderThroughNet, placedMessageReceived,
				networkSendCh, networkReceive, orderToSelf)
		}
	}

}
func InformationSharingThroughNet(ownElevator <- chan system.Elevator, broadcastElevatorInfo chan <- system.ElevatorInformation, receiveElevatorInfo <- chan system.ElevatorInformation,
	elevatorInfoCh chan<- system.ElevatorInformation) {
	for {
		select {
		case ownElev := <-ownElevator:
			//fmt.Printf("Elevatorinfo broadcasted: %#v\n", system.ElevatorInformation{ID: system.ElevatorID, Floor: ownElev.Floor, MotorDirection: ownElev.MotorDirection, Orders: ownElev.Orders, Behaviour: ownElev.Behaviour})
			broadcastElevatorInfo <- system.ElevatorInformation{ID: system.ElevatorID, Floor: ownElev.Floor, MotorDirection: ownElev.MotorDirection, Orders: ownElev.Orders, Behaviour: ownElev.Behaviour}

		case rcvElevInfo := <-receiveElevatorInfo:
			if rcvElevInfo.ID != system.ElevatorID {
				//fmt.Printf("Elevatorinfo from other elevator: %#v\n", rcvElevInfo)
				elevatorInfoCh <- rcvElevInfo
			}

		}
	}
}

func placeOrderNetworking(threadElevatorID int, sendingOrderThroughNet <-chan system.SendingOrder, placedMessageRecieved chan<- system.SendingOrder, networkSend chan<- system.SendingOrder, networkReceive <-chan system.SendingOrder, orderToSelf chan<- system.ButtonEvent) {
	for {
		select {
		case sOrdNet := <-sendingOrderThroughNet:
			if sOrdNet.ReceivingElevatorID == threadElevatorID {
				fmt.Printf("Order sent through network: %#v\n", sOrdNet)
				for i := 0; i < resendNum; i++ {
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
				for i := 0; i < resendNum; i++ {
					networkSend <- netReceive //As placed message }
				}
			}
		}
	}
}

func GetPeers(receivePeers <-chan peers.PeerUpdate, elevatorIDConnected chan <- int, elevatorIDDisconnected chan <- int) { //bør ID-ene være int eller string?
	for {
		select {
		case recPeer := <-receivePeers:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", recPeer.Peers)
			fmt.Printf("  New:      %q\n", recPeer.New)
			fmt.Printf("  Lost:     %q\n", recPeer.Lost)

			//hvis jeg får en ny bestilling OG den ikke er vår egen heis skal noe printes
			if recPeer.New != "" && recPeer.New != strconv.Itoa(system.ElevatorID) {
				fmt.Println("New peer ID: " + recPeer.New)
				newSentID, _ := strconv.Atoi(recPeer.New)
				elevatorIDConnected <- newSentID
			}
			for IDLost := 0; IDLost < len(recPeer.Lost); IDLost ++{
				lostSentID,_ := strconv.Atoi(recPeer.Lost[IDLost])
				fmt.Println("Lost sent ID:", lostSentID)
				elevatorIDDisconnected <- lostSentID
			}

		}
	}
}