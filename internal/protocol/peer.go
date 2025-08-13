package protocol

import "net"

type Peer interface {
	net.Conn
}
