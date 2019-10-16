package smf_message

var SmfChannel chan HandlerMessage

const (
	MaxChannel int = 1000
)

func init() {
	SmfChannel = make(chan HandlerMessage, MaxChannel)
}

func SendMessage(msg HandlerMessage) {
	SmfChannel <- msg
}
