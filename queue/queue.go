package queue

import "github.com/tetriscode/commander/proto"

//Producer is an abstraction over the queue impl
type Producer interface {
	NewProducer()
	SendCommand(*proto.Command)
	SendEvent(*proto.Event)
}

//Consumer is an abstraction over the queue impl
type Consumer interface {
	StartConsumer([]string, chan interface{})
	StopConsumer()
}

type Queue struct {
	Producer Producer
	Consumer Consumer
}

func NewQueue(c Consumer, p Producer) *Queue {
	return &Queue{p, c}
}
