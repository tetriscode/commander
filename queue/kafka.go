package queue

import "github.com/confluentinc/confluent-kafka-go/kafka"

type kafkaConsumer struct {
	c *kafka.Consumer
	topic string
}

type kafkaProducer struct {
	p *kafka.Producer
	topic string
}

func NewKafkaConsumer(topic string) Consumer {
	c := kafka.NewConsumer(&kafka.COnfigMap{
		"bootstrap.servers":               config.BootstrapServers,
		"group.id":                        groupID,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"}
	})

	return &kafkaConsumer{c,topic}
}

func (k *kafkaConsumer) StartConsumer() {
	err := k.c.SubscribeTopics([]string[k.topic],nil)
	if err != nil {
		log.Fatal("Failed to subscribe to topic:%s\n", k.topic)
	}
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	running := true
	for running == true {
		select {
		case sig := <-sigchan:
			utils.Logger.Printf("Caught signal %v: terminating kafka consumer: %s on: %s\n", sig, k.c, k.topic)
			running = false
		case evt := <-k.c.Events():
			switch e := evt.(type) {
			case kafka.AssignedPartitions:
				k.c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				k.c.Unassign()
			case *kafka.Message:
				receivedChan <- e
			case kafka.Error:
				log.Printf("%% Error: %v\n", e)
				running = false
			}
		}
	}
}

func (k *kafkaConsumer) StopConsumer() {
}

func (k *kafkaProducer) StartProducer() {

}

func (k *kafkaProducer) StopProducer() {

}
