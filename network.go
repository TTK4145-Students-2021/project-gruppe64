package main // dette må fikses
import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

	//"os"

)


type Elevator struct {
	ID     int
	Floor int
	MotorDirection int 		//egentlig hardwareIO.MotorDirection
	Orders [3][4] int 		// egentlig [hardwareIO.NumFloors][hardwareIO.NumButtons] int
	Behaviour string 		 //egentlig ElevatorBehaviour
	//tenker og at vi kanskje kan drite i å sende configStructen
}

//serializeElevator takes in an Elevator type struct, and serializes it so that it can be sent to another network.
func serializeElevator(elevatorStruct Elevator) []byte {
	serializedStruct, err := json.Marshal(elevatorStruct)
	if err != nil {
		log.Fatal(err)
	}
	return serializedStruct
}

//serializeOrder takes in a slice of strings, and serializes it so that it can be sent to another network.
func serializeOrder(orderSlice []interface{}) []byte{ //[]string) []byte{
	serializedSlice, err := json.Marshal(orderSlice)
	if err != nil {
		log.Fatal(err)
	}
	return serializedSlice
}

//deserializeElevator takes in a Serialized struct called Elevator, and deserializes it so that it can be used by the
//Order distributor module and the fsm module.
func deserializeElevator(serializedElevator []byte) Elevator{
	deserializedElevator := Elevator{}
	error := json.Unmarshal(serializedElevator, &deserializedElevator)
	if error != nil {
		log.Fatal(error)
	}
	return deserializedElevator
}

//deserializeOrder takes in a Serialized order, and deserializes it into a slice so that it can be used by the
//Order distributor module and the fsm module.
func deserializeOrder(serializedOrder []byte) []interface{}{
	deserializedOrder := []interface{}{[]int{1}, []int {1,2}}
	error := json.Unmarshal(serializedOrder, &deserializedOrder)
	if error != nil {
		log.Fatal(error)
	}
	return deserializedOrder
}


//broadcastElevatorStruct gets the Elevator type struct from the fsm module and sends it to a channel all elevators are
//listening to using UDP.
func broadcastElevatorStruct(sIP string, port int, ownElevatorStruct Elevator) {
	UDPAddr := net.UDPAddr{
		Port: port,
		IP: net.ParseIP(sIP), //mulig dette bør være en parse! Nei, da funker det ikke
	}

	conn, err := net.DialUDP("udp", nil, &UDPAddr) //er dette dumt å ha dette med dersom noden dør?
	if err != nil {
		log.Fatal(err)
	}

	for {
		ownElevator := serializeElevator(ownElevatorStruct)
		//message := []byte(strconv.Itoa(count))
		fmt.Println("Sent to server through UDP: ", string(ownElevator))//should be removed sooner or later
		_, err = conn.Write(ownElevator)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}

//Slår disse to sammen at then moment hvertfall
//placedOrderConfirmation takes NewAssignedOrder from the network module of another elevator as input, and sends an
//acknowledgement, OrderPlaced back.
//sendOrder receives an elevatorOrder from the order distributor module and sends it to the designated port of the
//elevator we want to take the order, using UDP.
func sendOrder(sIP string, port int,  orderToSend []interface{}){//[]string) {
	UDPAddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(sIP), //mulig dette bør være en parse!
	}

	conn, err := net.DialUDP("udp", nil, &UDPAddr) //er dette dumt å ha dette med dersom noden dør? Nei
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++{ //spam 10 messages
		serializedOrder := serializeOrder(orderToSend)
		fmt.Println("Sent to server through UDP on port", port,": ", string(serializedOrder)) //strconv.Itoa(count))
		_, err = conn.Write(serializedOrder)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond)//should change this to e.g time.Millisecond*10 in the actual code
		
	}
}



//readElevatorBroadcast takes in the elevatorstructs the other elevators send on the listening port, deserialize them
//and send them to the order designator.
func readElevatorBroadcast(port int) {//Elevator{
	UDPAddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(""),
	}
	conn, err := net.ListenUDP("udp", &UDPAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		elevatorStruct := make([]byte, 1024)
		n, UDPClient, err := (*conn).ReadFromUDP(elevatorStruct)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(UDPClient, " said on port ", port, ": ", string(elevatorStruct[:n]))
		//deserializedStruct := deserializeElevator(elevatorStruct[:n])
		//return deserializedStruct
	}
}


//Slår disse to sammen at the moment hvertfall:
//takeAssignedorder takes NewAssignedOrder from the networkmodule of another elevator as input, and sends it to the
//order distributor module after having deserialized it.
//confirmOrder will take in the acknowledgement-message from the networkmodule of another elevator and send
//that the order has been placed to the order distributor module.
func confirmOrder(port int) []interface{}{
	UDPAddr := net.UDPAddr{
	Port: port, //dette bør være heisens designerte port
	IP:   net.ParseIP(""), //mulig dette og bør være heisens designerte IP-adresse, men det er kanskje ikke så farlig
	}
	conn, err := net.ListenUDP("udp", &UDPAddr)
 	if err != nil {
		log.Fatal(err)
		}

	for {
		order := make([]byte, 1024)
		n, UDPClient, err := (*conn).ReadFromUDP(order)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(UDPClient, " said on port ", port, ": ", string(order[:n]))
		//deserializedOrder := deserializeOrder(order[:n]) //må bruke [:n] her for å kunne lese av json-objektet
		//return deserializedOrder
	}
}

/*test i main

func main() {
	comPort := 20001
	elev1Port := 20002
	//elev2Port := 20003
	//elev3Port := 20004
	elev1IP := "10.22.9.67"
	//elev2IP := "10.22.9.67"
	//elev3IP := "10.22.9.67"



	testSlice := []interface{}{[]int{1}, []int {1,2}}
	fmt.Println(reflect.TypeOf(testSlice)) //finne ut hva testSlicen er

	testSlice[0] = 1
	testSlice[1] = []int{1,2}


	elevator := Elevator{
		ID:     1,
		Behaviour: "Idle", //resten initialiseres til 0 eller tomme strenger
	}
	elevator.Orders[1][1] = 1
	elevator.Orders[0][0] = 1
	serEl := serializeElevator(elevator)
	//os.Stdout.Write(serEl) // en annen måte å skrive det ut, usikker på hva som er best
	fmt.Println(string(serEl))

	serSlice := serializeOrder(testSlice)
	fmt.Println(string(serSlice))
	a:=deserializeOrder(serSlice)
	fmt.Println(a[0])
	fmt.Println(a[1])
	go readElevatorBroadcast(comPort)
	go sendOrder(elev1IP, elev1Port, testSlice)
	go broadcastElevatorStruct(elev1IP, comPort, elevator) //funker å sende
	go confirmOrder(elev1Port)//skal egentlig være elev1Port her obvs
	select{}

//Broadcast and reading elevator at port 20001 works really well, order send and read does not work.


}

 */