package queue

import "github.com/tetriscode/commander/model"

//Producer is an abstraction over the queue impl
type Producer interface {
	SendCommand(*model.CommandParams) (*model.Command, error)
}

//Consumer is an abstraction over the queue impl
type Consumer interface {
	StartConsumer(func(v interface{}) error) error
	StopConsumer()
}

type Queue struct {
	Producer Producer
	Consumer Consumer
}

func NewQueue(c Consumer, p Producer) *Queue {
	return &Queue{p, c}
}
