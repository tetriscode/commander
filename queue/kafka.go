package queue

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/signal-go/utils"
)

type kafkaConsumer struct {
	c         *kafka.Consumer
	isRunning bool
	topic     string
}

type kafkaProducer struct {
	p     *kafka.Producer
	topic string
}

func NewKafkaConsumer(topic string) Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               "localhost:9092",
		"group.id":                        "commander",
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"}})
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil
	}
	return &kafkaConsumer{c, false, topic}
}

func (k *kafkaConsumer) StartConsumer(receivedChan chan *kafka.Message) {
	err := k.c.SubscribeTopics([]string{k.topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic:%s\n", k.topic)
	}
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	k.isRunning = true
	for k.isRunning == true {
		select {
		case sig := <-sigchan:
			utils.Logger.Printf("Caught signal %v: terminating kafka consumer: %s on: %s\n", sig, k.c, k.topic)
			k.isRunning = false
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
				k.isRunning = false
			}
		}
	}
	k.c.Close()
}

func (k *kafkaConsumer) StopConsumer() {
	if k.isRunning {
		k.isRunning = false
	}
}

func NewKafkaProducer() kafkaProducer {
	return kafkaProducer{kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"})}
}

func (k *kafkaProducer) SendCommand(cmd *model.Command) error {
	msg, err := k.sendMessage(k.topic, []byte(key), cmd...)
	if err != nil {
		return err
	}
	cmd.Timestamp = msg.Timestamp
	cmd.Topic = msg.TopicPartition.Topic
	cmd.Partition = msg.TopicPartition.Partition
	cmd.Offset = msg.TopicPartition.Offset
	cmd.Id = string(msg.Key)
	return nil
}

func (k *kafkaProducer) SendEvent(evt *model.Event) error {
	msg, err := k.sendMessage(k.topic, []byte(key), evt...)
	if err != nil {
		return err
	}
	evt.Timestamp = msg.Timestamp
	evt.Topic = msg.TopicPartition.Topic
	evt.Partition = msg.TopicPartition.Partition
	evt.Offset = msg.TopicPartition.Offset
	evt.Id = string(msg.Key)
	return nil
}

func (k *kafkaProducer) sendMessage(topic string, key, value []byte) (*kafka.Message, error) {
	deliveryChan := make(chan kafka.Event)
	err := k.p.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{
		Topic:     &topic,
		Partition: kafka.PartitionAny}, Key: key, Value: value},
		deliveryChan)

	if err != nil {
		return nil, err
	}

	del := <-deliveryChan
	msg := del.(*kafka.Message)

	defer close(deliveryChan)

	if msg.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", msg.TopicPartition.Error)
		return nil, msg.TopicPartition.Error
	}

	return msg, nil
}
