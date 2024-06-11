package kafka_api

import (
	"context"
	"fmt"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"strconv"
	"time"
)

// confluent-kafka-go配置参考librdkafka的配置https://docs.confluent.io/platform/current/clients/librdkafka/html/md_CONFIGURATION.html

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
		"client.id":                 "my_test_confluent_kafka_go_client",
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
	// 自动提交Offset的消费
	consumer := InitConsumer()
	defer consumer.Close()
	err := consumer.SubscribeTopics([]string{"qjj"}, nil)
	if err != nil {
		return err
	}

	cnt := 0

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			cnt++
			fmt.Printf("[%d]Message on %s: %s\n", cnt, msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
	return nil
}

//func ManuallyOffsetConsumer() error {
//	// 手动提交Offset的消费
//	cli := InitKafkaAdminClient()
//	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
//		"bootstrap.servers":         "localhost:9092",
//		"api.version.request":       "true",
//		"auto.offset.reset":         "latest",
//		"client.id":                 "confluent_kafka_go_client",
//		"heartbeat.interval.ms":     3000,
//		"session.timeout.ms":        30000,
//		"max.poll.interval.ms":      120000,
//		"fetch.max.bytes":           1024000,
//		"max.partition.fetch.bytes": 256000,
//		"enable.auto.commit":        false,
//		"group.id":                  "confluent-kafka-group"})
//	if err != nil {
//		return err
//	}
//	defer consumer.Close()
//
//	// 订阅的Topics列表
//	topics := []string{"qjj"}
//	err = consumer.SubscribeTopics(topics, nil)
//	if err != nil {
//		return err
//	}
//
//	// 获取当前的offsets
//	for _, topic := range topics {
//
//		offsets, err := cli.ListConsumerGroupOffsets(group, []kafka.TopicPartition{
//			{Topic: &topic, Partition: kafka.PartitionAny},
//		})
//		if err != nil {
//			panic(err)
//		}
//	}
//
//
//	cnt := 0
//	offsetMap := map[string]map[int32]kafka.Offset{}
//
//	go func() {
//		for {
//			msg, err := consumer.ReadMessage(-1)
//			if err == nil {
//				cnt++
//				fmt.Printf("[%d]Message on %s: %s\n", cnt, msg.TopicPartition, string(msg.Value))
//
//				tp := msg.TopicPartition
//				topic := *(tp.Topic)
//				id := tp.Partition
//				offset := tp.Offset
//				_, exists := offsetMap[topic]
//				if exists {
//					i, exists := offsetMap[topic][id]
//					if (exists && offset > i) || !exists {
//						offsetMap[topic][id] = offset
//					}
//				} else {
//					topicOffsetMap := map[int32]kafka.Offset{}
//					topicOffsetMap[id] = offset
//					offsetMap[topic] = topicOffsetMap
//				}
//
//			} else {
//				// The client will automatically try to recover from all errors.
//				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
//			}
//
//		}
//	}()
//
//	for {
//		if cnt%1000 == 0 {
//			for t, oi := range offsetMap {
//				var l []kafka.TopicPartition
//				for p, o := range oi {
//					l = append(l, kafka.TopicPartition{Topic: &t, Partition: p, Offset: o})
//				}
//				fmt.Printf("[topic:%s cnt:%d]Commit. \n", t, cnt)
//				_, err := consumer.CommitOffsets(l)
//				if err != nil {
//					println(err.Error())
//				}
//			}
//		}
//	}
//	return nil
//}

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

func InitKafkaAdminClient() *kafka.AdminClient {
	cli, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		panic(err)
	}
	fmt.Print("init kafka admin client success\n")
	return cli
}

func UseKafkaAdminClient() {
	ac := InitKafkaAdminClient()
	topic := "qjj"
	group := "confluent-kafka-group"

	var groupsPartitions []kafka.ConsumerGroupTopicPartitions
	groupsPartitions = append(groupsPartitions, kafka.ConsumerGroupTopicPartitions{
		Group: group,
		Partitions: []kafka.TopicPartition{
			kafka.TopicPartition{
				Topic:     &topic,
				Partition: 0,
			},
			kafka.TopicPartition{
				Topic:     &topic,
				Partition: 1,
			},
		},
	})
	res, err := ac.ListConsumerGroupOffsets(context.Background(), groupsPartitions)
	if err != nil {
		panic(err)
	}

	offsetMap := map[string]map[int32]kafka.Offset{}
	offsetMap[topic] = map[int32]kafka.Offset{}

	for _, partition := range res.ConsumerGroupsTopicPartitions {
		partitions := partition.Partitions
		for _, topicPartition := range partitions {
			p := topicPartition.Partition
			o := topicPartition.Offset
			offsetMap[topic][p] = o
		}
	}

	println(fmt.Sprintf("offsetMap: %v", offsetMap))
}
