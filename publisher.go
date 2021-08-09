package walrus

type publisher struct {
	dispatcher Dispatcher
}

type Publisher interface {
	Publish(message []byte)
}

func NewPublisher(d Dispatcher) Publisher {
	return publisher{
		dispatcher: d,
	}
}

func (p publisher) Publish(message []byte) {
	p.dispatcher.Dispatch(message)
}
