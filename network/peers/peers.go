package peers

import (
	"fmt"
	"net"
	"realtimeProject/project-gruppe64/network/conn"
	"sort"
	"time"
)

type PeerUpdate struct {
	Peers []string
	New   string //trenger vi egentlig new og lost? Usikker
//	Lost  []string
}

const interval = 15 * time.Millisecond
const timeout = 500 * time.Millisecond

//Transmitter får igang en samtale på en port og overfører etter en viss tid signalene sine dit. Dette kan brukes som en
//"jeg er ikke død-greie om man vil. Tror ikke vi trenger det, men hvem vet
func Transmitter(port int, id string, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port)) //denne bør kanskje endres

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)//[]byte(id), addr)
		}
	}
}

//Receiver mottar peerUpdaten på en port, og leser av nye, døde og gamle noder
func Receiver(port int, peerUpdateCh chan<- PeerUpdate) {

	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// Adding new connection
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				//p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}


		// Removing dead connection
		/*Tror ikke denne trengs heller
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}
		*/


		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			//sort.Strings(p.Lost)
			peerUpdateCh <- p
		}
	}
}
