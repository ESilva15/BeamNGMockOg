package mockserver

import (
	"fmt"
	"net"
)

func openUDPServer(addr string, port int) (*net.UDPAddr, *net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil, nil, err
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, nil, err
	}

	return udpAddr, conn, nil
}
