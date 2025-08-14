package protocol

const (
	EXIT_OK                    uint8 = 0x00
	EXIT_CLIENT_ALREADY_EXISTS uint8 = 0x01
	EXIT_NO_SUCH_CLIENT        uint8 = 0x02
	EXIT_NO_NEW_MESSAGES       uint8 = 0x03
)
