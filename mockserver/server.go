package mockserver

import (
	"fmt"
	"net"
)

type Transport interface {
	Send(data []byte) error
	Close() error
}

type UDPTransport struct {
	Conn *net.UDPConn
	Addr *net.UDPAddr
}

func NewUDPTransport(address string, port int) (*UDPTransport, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, err
	}

	return &UDPTransport{Conn: conn, Addr: addr}, nil
}

// Send will send a byte array of data trough the UDP server
func (u *UDPTransport) Send(data []byte) (int, error) {
	return u.Conn.WriteToUDP(data, u.Addr)
}

func (u *UDPTransport) Close() error {
	return u.Conn.Close()
}
