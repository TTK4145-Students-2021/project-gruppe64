package sendandreceive

import (
	"fmt"
	"realtimeProject/project-gruppe64/network/bcast"
	"realtimeProject/project-gruppe64/system"
)

const (
	resendNum = 10
	)



func SetUpReceiverAndTransmitterPorts(receiveElevatorInfo chan system.ElevatorInformation, broadcastElevatorInfo chan system.ElevatorInformation,
	networkReceive chan system.SendingOrder, sendingOrderThroughNet <-chan system.SendingOrder,
	messageTimer chan<- system.SendingOrder, orderToSelf chan<- system.ButtonEvent){

	go bcast.Receiver(60000, receiveElevatorInfo) //Receive others elevator information
	go bcast.Transmitter(60000, broadcastElevatorInfo) //Send elevator Information

	 //PLACED MESSAGE AND ORDER ON SAME. TO FROM/
	go bcast.Receiver(60001+system.ElevatorID, networkReceive) //Receive orders

	for elevID := 0; elevID < system.NumElevators; elevID++ {
		if elevID != system.ElevatorID {
			networkSendCh := make(chan system.SendingOrder) //Reset every run
			go bcast.Transmitter(60001 +elevID, networkSendCh) //Transmit orders to place
			go placeOrderNetworking(elevID, sendingOrderThroughNet, messageTimer,
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

func placeOrderNetworking(threadElevatorID int, sendingOrderThroughNet <-chan system.SendingOrder, messageTimer chan<- system.SendingOrder, networkSend chan<- system.SendingOrder, networkReceive <-chan system.SendingOrder, orderToSelf chan<- system.ButtonEvent) {
	duplicate := system.SendingOrder{}
	for{
		select{
		case sOrdNet := <-sendingOrderThroughNet:
			if sOrdNet.ReceivingElevatorID == threadElevatorID {
				fmt.Printf("Order sent through network: %#v\n", sOrdNet)
				for i := 0; i < resendNum; i++ {
					//time.Sleep(1 * time.Millisecond)
					networkSend <- sOrdNet
				}
			}
		case netReceive := <-networkReceive:
			if netReceive.SendingElevatorID == system.ElevatorID { //THEN IT IS A PLACED MESSAGE
				if duplicate != netReceive {
					fmt.Println("Placed message reveived")
					messageTimer <- netReceive
					duplicate = netReceive
				}
			}

			if netReceive.ReceivingElevatorID == system.ElevatorID { //THEN IT IS A ORDER
				if duplicate != netReceive {
					fmt.Printf("Order received: %#v\n", netReceive)
					orderToSelf <- netReceive.Order
					for i := 0; i < resendNum; i ++ {
						networkSend <- netReceive //As placed message
					}
				}
			}
		}
	}
}
