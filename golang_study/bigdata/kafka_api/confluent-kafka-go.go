package kafka_api

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"strconv"
	"time"
)

const (
	INT32_MAX = 2147483647 - 1000
)

type KafkaConfig struct {
	Topics           []string `json:"topics"`
	GroupId          string   `json:"group.id"`
	BootstrapServers string   `json:"bootstrap.servers"`
	SecurityProtocol string   `json:"security.protocol"`
	SaslMechanism    string   `json:"sasl.mechanism"`
	SaslUsername     string   `json:"sasl.username"`
	SaslPassword     string   `json:"sasl.password"`
	AutoOffsetReset  string   `json:"auto.offset.reset"`
}

func InitConsumer() *kafka.Consumer {
	var kafkaconf = &kafka.ConfigMap{
		"api.version.request":       "true",
		"auto.offset.reset":         "latest",
		"heartbeat.interval.ms":     3000,
		"session.timeout.ms":        30000,
		"max.poll.interval.ms":      120000,
		"fetch.max.bytes":           1024000,
		"max.partition.fetch.bytes": 256000}
	kafkaconf.SetKey("bootstrap.servers", "localhost:9092")
	kafkaconf.SetKey("group.id", "confluent-kafka-go")

	consumer, err := kafka.NewConsumer(kafkaconf)
	if err != nil {
		panic(err)
	}
	fmt.Print("init kafka consumer success\n")
	return consumer
}

func DoConsume() error {
	consumer := InitConsumer()
	defer consumer.Close()
	err := consumer.SubscribeTopics([]string{"qjj"}, nil)
	if err != nil {
		return err
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will
			//automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
	return nil
}

func InitProducer() *kafka.Producer {
	var kafkaconf = &kafka.ConfigMap{
		"api.version.request":           "true",
		"message.max.bytes":             1000000,
		"linger.ms":                     500,
		"sticky.partitioning.linger.ms": 1000,
		"retries":                       INT32_MAX,
		"retry.backoff.ms":              1000,
		"acks":                          "1"}

	kafkaconf.SetKey("bootstrap.servers", "localhost:9092")
	producer, err := kafka.NewProducer(kafkaconf)
	if err != nil {
		panic(err)
	}
	fmt.Print("init kafka producer success\n")
	return producer
}

func DoProduce() {
	producer := InitProducer()
	defer producer.Close()
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed to write access log entry:%v", ev.TopicPartition.Error)
				} else {
					log.Printf("Send OK topic:%v partition:%v offset:%v content:%s\n", *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset, ev.Value)

				}
			}
		}
	}()

	Topic1 := "qjj"
	Topic2 := "t_20111_3"
	// Produce messages to topic (asynchronously)
	i := 0
	for {
		i = i + 1
		value := "this is a kafka message from confluent go " + strconv.Itoa(i)
		var msg *kafka.Message = nil
		if i%2 == 0 {
			msg = &kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &Topic1, Partition: kafka.PartitionAny},
				Value:          []byte(value),
			}
		} else {
			msg = &kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &Topic2, Partition: kafka.PartitionAny},
				Value:          []byte(value),
			}
		}
		producer.Produce(msg, nil)
		time.Sleep(time.Duration(1) * time.Millisecond)
	}
	// Wait for message deliveries before shutting down
	producer.Flush(15 * 1000)

}
