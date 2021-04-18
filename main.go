package main

import (
	"fmt"
	"io/ioutil"
	"realtimeProject/project-gruppe64/distributor"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/network/peers"
	"realtimeProject/project-gruppe64/network/sendandreceive"
	"realtimeProject/project-gruppe64/system"
	"realtimeProject/project-gruppe64/timer"
	"runtime"
	"strconv"
	"time"
)

func primaryWork(activateAsPrimary <-chan bool){
	for{
		select{
		case activate := <-activateAsPrimary:
			if activate{
				fmt.Println(activate)
				hardwareIO.Init(system.LocalHost, system.NumFloors)
				system.SpawnBackup()


				orderToSelfCh := make(chan system.ButtonEvent)
				hallOrderCh := make(chan system.ButtonEvent)
				floorArrivalCh := make(chan int)
				obstructionEventCh := make(chan bool)
				doorTimerDurationCh := make(chan float64)
				doorTimerTimedOutCh := make(chan bool)

				ownElevatorCh := make(chan system.Elevator)

				//Network channels for transmitting and receiving
				receiveElevatorInfoCh := make(chan system.ElevatorInformation)
				broadcastElevatorInfoCh := make(chan system.ElevatorInformation)
				networkReceiveCh := make(chan system.SendingOrder)
				elevatorIDConnectedCh := make(chan int) //DENNE HG, sender id til heis når connected
				elevatorIDDisconnectedCh := make(chan int) //DENNE HG, sender id til heis når disconnected
				receivePeersCh := make(chan peers.PeerUpdate)


				sendingOrderThroughNetCh := make(chan system.SendingOrder) //channel that receives
				elevatorInfoCh := make(chan system.ElevatorInformation) //channel with elevatorinformation, sent from networkmodule to
				//own modules

				messageTimerCh := make(chan system.SendingOrder)
				placedMessageReceivedCh := make(chan system.SendingOrder)
				messageTimerTimedOutCh := make(chan system.SendingOrder)
				orderTimerCh := make(chan system.SendingOrder)
				orderTimerTimedOutCh := make(chan system.SendingOrder)

				// Timers:
				go timer.RunDoorTimer(doorTimerDurationCh, doorTimerTimedOutCh)
				go timer.RunMessageTimer(messageTimerCh, placedMessageReceivedCh, messageTimerTimedOutCh)
				go timer.RunOrderTimer(orderTimerCh, orderTimerTimedOutCh)

				// Hardware:
				go hardwareIO.RunHardware(orderToSelfCh, hallOrderCh, floorArrivalCh, obstructionEventCh)

				// Distributor and FSM:
				go fsm.ElevatorFSM(orderToSelfCh, floorArrivalCh, obstructionEventCh, ownElevatorCh, doorTimerDurationCh, doorTimerTimedOutCh)
				go distributor.OrderDistributor(hallOrderCh, elevatorInfoCh, ownElevatorCh, sendingOrderThroughNetCh, orderToSelfCh, messageTimerCh, messageTimerTimedOutCh, orderTimerCh, orderTimerTimedOutCh, elevatorIDConnectedCh, elevatorIDDisconnectedCh)

				//Network:
				go sendandreceive.SetUpReceiverAndTransmitterPorts(receiveElevatorInfoCh, broadcastElevatorInfoCh, networkReceiveCh, sendingOrderThroughNetCh, placedMessageReceivedCh, orderToSelfCh, receivePeersCh)
				go sendandreceive.InformationSharingThroughNet(ownElevatorCh, broadcastElevatorInfoCh, receiveElevatorInfoCh, elevatorInfoCh)
				go sendandreceive.GetPeers(receivePeersCh, elevatorIDConnectedCh, elevatorIDDisconnectedCh)

				docNum := 0
				for {
					_ = ioutil.WriteFile("system/primary_doc.txt", []byte(strconv.FormatInt(int64(docNum), 10)), 0644)
					time.Sleep(1*time.Second)
					docNum += 1
				}
			}
		default:
			break
		}
	}
}




func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	activateAsPrimaryCh := make(chan bool)

	if system.IsBackup(){
		go system.CheckPrimaryExistence(activateAsPrimaryCh)
		go primaryWork(activateAsPrimaryCh)
	} else {
		system.MakeBackupFile()
		go primaryWork(activateAsPrimaryCh)
		activateAsPrimaryCh <- true
	}


	for {}
}