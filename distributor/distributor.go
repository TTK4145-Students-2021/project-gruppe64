package distributor

import (
	"fmt"
	"realtimeProject/project-gruppe64/configuration"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/network/sendandreceive"
	"strconv"
)





// GOROUTINE:
func OrderDistributor(hallOrder <-chan hardwareIO.ButtonEvent, elevatorInfo <-chan sendandreceive.ElevatorInformation, ownElevator <-chan fsm.Elevator, sendingOrderThroughNet chan<- sendandreceive.SendingOrder, orderToSelf chan<- hardwareIO.ButtonEvent, messageTimer chan<- sendandreceive.SendingOrder, messageTimerTimedOut <-chan sendandreceive.SendingOrder, orderTimer chan<- sendandreceive.SendingOrder, orderTimerTimedOut <- chan sendandreceive.SendingOrder){
	elevs := initiateElevatorsTagged()
	var distributedOrders map[string][]sendandreceive.SendingOrder
	for {
		select {
		case hallOrd := <-hallOrder:
			switch hallOrd.Button {
			case hardwareIO.BT_HallUp:
				elevs.HallOrders[hallOrd.Floor][0] = true
			case hardwareIO.BT_HallDown:
				elevs.HallOrders[hallOrd.Floor][1] = true
			default:
				break
			}
			designatedID := getDesignatedElevatorID(elevs)
			if designatedID == configuration.ElevatorID {
				fmt.Println("ORDER TO SELF")
				orderToSelf <- hallOrd
			} else {
				sOrd := sendandreceive.SendingOrder{designatedID, configuration.ElevatorID, hallOrd}
				sendingOrderThroughNet <- sOrd
				messageTimer <- sOrd
				orderTimer <- sOrd
				distributedOrders[strconv.Itoa(designatedID)] = append(distributedOrders[strconv.Itoa(designatedID)], sOrd)
			}
			switch hallOrd.Button { //sletter ordren fra her og nå
			case hardwareIO.BT_HallUp:
				elevs.HallOrders[hallOrd.Floor][0] = false
			case hardwareIO.BT_HallDown:
				elevs.HallOrders[hallOrd.Floor][1] = false
			default:
				break
			}

		case msgTimedOut := <- messageTimerTimedOut:
			orderToSelf <- msgTimedOut.Order

		case ordTimedOut := <- orderTimerTimedOut:
			for key, dOrds := range distributedOrders{
				if key == strconv.Itoa(ordTimedOut.ReceivingElevatorID){
					for _, dOrd := range dOrds{
						if dOrd == ordTimedOut{
							orderToSelf <- ordTimedOut.Order
						}
					}
				}
			}
		case elev := <-ownElevator:
			elevs.States[strconv.Itoa(configuration.ElevatorID)] = getUpdatedElevatorTagged(sendandreceive.ElevatorInformation{configuration.ElevatorID, elev.Floor, elev.MotorDirection, elev.Orders, elev.Behaviour})

		case elevInfo := <-elevatorInfo:
			if removeExecutedOrders(elevInfo, distributedOrders[strconv.Itoa(elevInfo.ID)]) != nil{ //måtte legge til dette for å ikke få feilmelding
				distributedOrders[strconv.Itoa(elevInfo.ID)] = removeExecutedOrders(elevInfo, distributedOrders[strconv.Itoa(elevInfo.ID)])
			}
			elevs.States[strconv.Itoa(elevInfo.ID)] = getUpdatedElevatorTagged(elevInfo)
		}
	}
}

func elevatorOperationCheck(){
	//Sjekk ordrer som er kommet inn sist, er de like de andre for så så lenge? I så fall opererer ikke heisen. Også en timer på når akk
	//den heisen mottok sist oppdatering. On lenge siden? tilsvarende som om de er like. Heisen opererer ikke.
	//Tenk litt som timerene for order og message
}


