package udp

import (
	"log"
	"net"
)

func CreateServer() {

	udpServer, err := net.ListenPacket("udp", ":3999")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	log.Println("UDP server started on port 3999")

	for {
		buf := make([]byte, 1024*16)
		_, addr, err := udpServer.ReadFrom(buf)

		if err != nil {
			log.Fatal(err)
			continue
		}

		udpServer.WriteTo([]byte(string(buf)), addr)
	}
}
