package sendandreceive

import (
	"fmt"
	"realtimeProject/project-gruppe64/network/bcast"
	"realtimeProject/project-gruppe64/system"
	"time"
)

const (
	resendNum = 1
)



func GetReceiverAndTransmitterPorts(receiveElevatorInfo chan system.ElevatorInformation, broadcastElevatorInfo chan system.ElevatorInformation,
	receiveOrder chan system.SendingOrder, sendPlacedMessage chan system.SendingOrder, sendingOrderThroughNet <-chan system.SendingOrder, messageTimer chan<- system.SendingOrder){

	go bcast.Receiver(60000, receiveElevatorInfo) //Receive others elevator information
	go bcast.Transmitter(60000, broadcastElevatorInfo) //Send elevator Information


	go bcast.Receiver(60001+system.ElevatorID, receiveOrder) //Receive orders
	go bcast.Transmitter(60002+system.ElevatorID, sendPlacedMessage) //Transmit placed message
	for elevID := 0; elevID < system.NumElevators; elevID++ {
		if elevID != system.ElevatorID {
			placeOrderCh := make(chan system.SendingOrder) //Reset every run
			receivePlacedMessageCh := make(chan system.SendingOrder)
			go bcast.Transmitter(60001 +elevID, placeOrderCh) //Transmit orders to place
			go bcast.Receiver(60002 +elevID, receivePlacedMessageCh) //Receive placed message
			go placeOrderNetworking(elevID, sendingOrderThroughNet, messageTimer,
				placeOrderCh, receivePlacedMessageCh)
		}
	}

}
func SendReceiveOrders(ownElevator <- chan system.Elevator, broadcastElevatorInfo chan <- system.ElevatorInformation, receiveElevatorInfo <- chan system.ElevatorInformation,
	elevatorInfoCh chan<- system.ElevatorInformation, receiveOrder <-chan system.SendingOrder, orderToSelf chan<- system.ButtonEvent, sendPlacedMessage chan<- system.SendingOrder) {
	for {
		select {
		case ownElev := <-ownElevator:
			//fmt.Printf("Elevatorinfo broadcasted: %#v\n", fsmElevatorToElevatorNetwork(e))
			broadcastElevatorInfo <- system.ElevatorInformation{ID: system.ElevatorID, Floor: ownElev.Floor, MotorDirection: ownElev.MotorDirection, Orders: ownElev.Orders, Behaviour: ownElev.Behaviour}

		case rcvElevInfo := <-receiveElevatorInfo:
			//fmt.Printf("Elevatorinfo from other elevators and own broadcasted: %#v\n", o)
			if rcvElevInfo.ID != system.ElevatorID { //FSM sender alt egen til distributor
				elevatorInfoCh <- rcvElevInfo
			}

		case rcvOrd := <-receiveOrder:
			fmt.Printf("Order received: %#v\n", rcvOrd)
			orderToSelf <- rcvOrd.Order
			sendPlacedMessage <- rcvOrd
		}
	}
}

func placeOrderNetworking(threadElevatorID int, sendingOrderThroughNet <-chan system.SendingOrder, messageTimer chan<- system.SendingOrder, placeOrder chan<- system.SendingOrder, receivePlacedMessage <-chan system.SendingOrder) {
	for{
		select{
		case sOrdNet := <-sendingOrderThroughNet:
			if sOrdNet.ReceivingElevatorID == threadElevatorID {
				fmt.Printf("Order sent through network: %#v\n", sOrdNet)
				for i := 0; i < resendNum; i++ {
					time.Sleep(1 * time.Millisecond)
					placeOrder <- sOrdNet
				}
			}
		case rcvPlcdMsg := <-receivePlacedMessage:
			fmt.Println("Placed message reveived")
			messageTimer <- rcvPlcdMsg
		}

	}
}
