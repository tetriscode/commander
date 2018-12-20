package queue

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/util"
)

type kafkaConsumer struct {
	c         *kafka.Consumer
	isRunning bool
	topics    []string
}

type kafkaProducer struct {
	p     *kafka.Producer
	topic string
}

var EVENTS_TOPIC, COMMANDS_TOPIC string

// NewKafkaConsumer creates a new kafka consumer
func NewKafkaConsumer(topics []string) (Consumer, error) {
	EVENTS_TOPIC = os.Getenv("KAFKA_EVENTS_TOPIC")
	COMMANDS_TOPIC = os.Getenv("KAFKA_COMMANDS_TOPIC")
	util.Log.Debug().Verb("creating").Object("kafka consumer").Log()
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               os.Getenv("KAFKA_BROKERS"),
		"group.id":                        os.Getenv("KAFKA_GROUP_ID"),
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kafka consumer")
	}
	util.Log.Debug().Verb("created").Object("kafka consumer").Log()
	return &kafkaConsumer{c, false, topics}, nil
}

func kafkaMessageToEntity(k *kafka.Message) interface{} {
	topic := k.TopicPartition.Topic
	if *topic == COMMANDS_TOPIC {
		util.Log.Debug().Verb("received").Object("command").IndirectObject("kafka").Log()
		var cmd model.CommandParams
		err := proto.Unmarshal(k.Value, &cmd)
		if err != nil {
			util.Log.Debug().Verb("errored parsing").Object("message").IndirectObject("kafka").Log()
			log.Fatal(err)
		} else {
			uid, err := uuid.Parse(string(k.Key))
			if err != nil {
				util.Log.Debug().Verb("errored parsing").Object("key").IndirectObject("message").Log()
				log.Fatal(err)
			}
			return &model.Command{Id: &model.UUID{Value: uid.String()},
				Action:    cmd.Action,
				Data:      cmd.Data,
				Topic:     *k.TopicPartition.Topic,
				Offset:    int64(k.TopicPartition.Offset),
				Timestamp: k.Timestamp.Unix(),
				//TODO Children
			}
		}
	} else if *topic == EVENTS_TOPIC {
		util.Log.Debug().Verb("received").Object("event").IndirectObject("kafka").Log()
		var evt model.Event
		err := proto.Unmarshal(k.Value, &evt)
		if err != nil {
			util.Log.Debug().Verb("errored parsing").Object("message").IndirectObject("kafka").Log()
			log.Fatal(err)
		} else {
			uid, err := uuid.Parse(string(k.Key))
			if err != nil {
				util.Log.Debug().Verb("errored parsing").Object("key").IndirectObject("message").Log()
				log.Fatal(err)
			}
			return &model.Event{Id: &model.UUID{Value: uid.String()},
				Action:    evt.Action,
				Data:      evt.Data,
				Topic:     *k.TopicPartition.Topic,
				Offset:    int64(k.TopicPartition.Offset),
				Timestamp: k.Timestamp.Unix(),
				Parent:    evt.GetParent(),
				//TODO Children
			}
		}
	}
	return nil
}

func (k *kafkaConsumer) StartConsumer(fn func(interface{}) error) error {
	util.Log.Debug().Verb("subscribes").Object(fmt.Sprint(k.topics)).IndirectObject("kafka").Log()
	err := k.c.SubscribeTopics(k.topics, nil)
	if err != nil {
		util.Log.Debug().Verb("errored subscribing").Object(fmt.Sprint(k.topics)).IndirectObject("kafka").Log()
		log.Fatalf("Failed to subscribe to topic:%s\n", k.topics)
	}
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		k.isRunning = true
		util.Log.Debug().Verb("starting").Object("loop").IndirectObject("kafka").Log()
		for k.isRunning == true {
			select {
			case sig := <-sigchan:
				util.Log.Debug().Verb("received").Object("sigterm").IndirectObject("kafka").Log()
				log.Printf("Caught signal %v: terminating kafka consumer: %s on: %s\n", sig, k.c, k.topics)
				k.isRunning = false
			case evt := <-k.c.Events():
				util.Log.Debug().Verb("received").Object("event").IndirectObject("kafka").Log()
				switch e := evt.(type) {
				case kafka.AssignedPartitions:
					k.c.Assign(e.Partitions)
				case kafka.RevokedPartitions:
					k.c.Unassign()
				case *kafka.Message:
					util.Log.Debug().Verb("received").Object("message").IndirectObject("kafka").Log()
					err := fn(kafkaMessageToEntity(e))
					if err != nil {
						log.Fatal(err.Error())
					}
				case kafka.Error:
					log.Printf("%% Error: %v\n", e)
					k.isRunning = false
				}
			}
		}
		util.Log.Debug().Verb("stopping").Object("loop").IndirectObject("kafka").Log()
		k.c.Close()
		util.Log.Debug().Verb("closed").Object("consumer").IndirectObject("kafka").Log()
	}()
	return nil
}

func (k *kafkaConsumer) StopConsumer() {
	if k.isRunning {
		k.isRunning = false
	}
}

// NewKafkaProducer creates a new kafka producer
func NewKafkaProducer(topic string) (Producer, error) {
	log.Println("creating new kafka producer")
	config := &kafka.ConfigMap{
		"bootstrap.servers":    os.Getenv("KAFKA_BROKERS"),
		"group.id":             os.Getenv("KAFKA_GROUP_ID"),
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kafka producer")
	}
	return &kafkaProducer{p: producer, topic: topic}, nil
}

func (k *kafkaProducer) SendCommand(cmdp *model.CommandParams) (*model.Command, error) {
	util.Log.Debug().Verb("sending").Object("message").IndirectObject("command").PrepObject("kafka").Log()
	pbf, err := proto.Marshal(cmdp)
	if err != nil {
		util.Log.Debug().Verb("errored sending").Object("message").IndirectObject("kafka").Log()
		return nil, err
	}
	id := uuid.New().String()
	msg, err := k.sendMessage(k.topic, []byte(id), pbf)
	if err != nil {
		return nil, err
	}
	util.Log.Debug().Verb("sent").Object("message").IndirectObject("command").PrepObject("kafka").Log()
	var cmd model.Command
	cmd.Timestamp = msg.Timestamp.Unix()
	cmd.Topic = *msg.TopicPartition.Topic
	cmd.Partition = msg.TopicPartition.Partition
	cmd.Offset = int64(msg.TopicPartition.Offset)
	cmd.Id = &model.UUID{Value: id}
	return &cmd, nil
}

func (k *kafkaProducer) sendMessage(topic string, key, value []byte) (*kafka.Message, error) {
	util.Log.Debug().Verb("sending").Object("message").IndirectObject(string(key)).PrepObject("kafka").Log()
	deliveryChan := make(chan kafka.Event)
	err := k.p.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{
		Topic:     &topic,
		Partition: kafka.PartitionAny}, Key: key, Value: value},
		deliveryChan)

	if err != nil {
		util.Log.Debug().Verb("errored sending").Object("message").IndirectObject(string(key)).PrepObject("kafka").Log()
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
