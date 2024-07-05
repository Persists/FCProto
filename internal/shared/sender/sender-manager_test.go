package sender

// SenderManager manages the sender
type MockedSenderManager struct {
	Sender   *Sender
	DataChan chan string
	StopChan chan bool
}
