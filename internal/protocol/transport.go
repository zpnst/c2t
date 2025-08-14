package protocol

// Transport interface for instance
type ITransport interface {
	Addr() string
	Close() error
	ListenAndAccept() error
	Consume() <-chan Message
	Peer(string) Peer
	SendAnswer(src string, b any) error
}

// Transport interface for client
type CTransport interface {
	DialInstance() error
	DecodeAnswer(*Message) error
	EncodeMessage(*Message) error
}
