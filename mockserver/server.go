package mockserver

import (
	"fmt"
	"net"
)

type udpServer struct {
	Addr *net.UDPAddr
	Conn *net.UDPConn
}

func newUDPServer(addr string, port int) (udpServer, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return udpServer{}, err
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return udpServer{}, err
	}

	return udpServer{Addr: udpAddr, Conn: conn}, nil
}

func (u *udpServer) Close() error {
	return u.Conn.Close()
}
