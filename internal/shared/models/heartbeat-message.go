package models

type HeartbeatMessage struct {
	CallbackPort string
}

func NewHeartbeatMessage(addr string) HeartbeatMessage {
	return HeartbeatMessage{
		CallbackPort: addr,
	}
}
