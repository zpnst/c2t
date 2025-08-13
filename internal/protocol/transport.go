package protocol

// Transport interface for instance
type ITransport interface {
	Addr() string
	Close() error
	ListenAndAccept() error
	Consume() <-chan Message
	Peer(string) Peer
}

// Transport interface for client
type CTransport interface {
	DialInstance() error
	Send(Message) error
	WaitForInstance() error
}
