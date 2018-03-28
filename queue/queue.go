package queue

type Producer interface {
	StartProducer()
	StopProducer()
}

type Consumer interface {
	StartConsumer()
	StopConsumer()
}

type Queue struct {
	producer Producer
	consumer Consumer
}

func NewQueue(c Consumer, p Producer) *Queue {

}
