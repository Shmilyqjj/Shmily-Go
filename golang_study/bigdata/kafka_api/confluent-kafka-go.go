package kafka_api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"strconv"
	"time"
)

// confluent-kafka-go配置参考librdkafka的配置https://docs.confluent.io/platform/current/clients/librdkafka/html/md_CONFIGURATION.html
// 由于confluent-kafka-go会调用C库librdkafka，编译时不能禁用CGO,当你的Go程序需要调用C库或依赖于C代码时，必须启用CGO（CGO_ENABLED=1），否则无法编译通过，报错undefined xxx

/**
CGO_ENABLED=1：启用 CGO。这允许 Go 程序调用 C 代码并链接 C 库。在编译过程中，如果你的 Go 代码中包含对 C 代码的引用，这个设置是必要的。
CGO_ENABLED=0：禁用 CGO。这意味着编译器将忽略所有的 C 代码和 C 库链接。这对于构建完全独立的二进制文件非常有用，因为禁用 CGO 后生成的二进制文件不需要任何外部 C 库的支持。
*/

const (
	INT32_MAX     = 2147483647 - 1000
	LineDelimiter = "\n"
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

func SampleManuallyOffsetConsumer() error {
	brokers := "localhost:9092"
	topics := []string{"qjj"}
	group := "confluent-kafka-group"

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	configMap := kafka.ConfigMap{
		"bootstrap.servers":         brokers,
		"api.version.request":       "true",
		"auto.offset.reset":         "latest",
		"broker.address.family":     "v4",
		"client.id":                 "confluent_kafka_go_client",
		"heartbeat.interval.ms":     3000,
		"session.timeout.ms":        30000,
		"max.poll.interval.ms":      120000,
		"fetch.max.bytes":           1024000,
		"max.partition.fetch.bytes": 256000,
		//"enable.auto.commit":       false, // 关闭自动提交偏移量（手动提交）
		"enable.auto.offset.store": false, // 关闭自动存储偏移量（手动存储）
		"group.id":                 group}

	consumer, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return err
	}
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return err
	}
	defer consumer.Close()

	// 消费消息
	totalCnt := 0
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(-1)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				// Process the message received.
				totalCnt++
				msg := fmt.Sprintf("[%d]Message on %s: %s", totalCnt, e.TopicPartition, string(e.Value))
				println(msg)
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
				_, err := consumer.StoreMessage(e)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%% Error storing offset after message %s:\n",
						e.TopicPartition)
				}
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
	return nil
}

func ManuallyOffsetConsumer() error {
	// 手动提交Offset的消费  批处理，达到一定条数或时间后进行批处理
	brokers := "localhost:9092"
	topics := []string{"qjj"}
	group := "confluent-kafka-group"
	batchInterval := 5 * time.Second
	maxBatchSize := 10000

	configMap := kafka.ConfigMap{
		"bootstrap.servers":         brokers,
		"api.version.request":       "true",
		"auto.offset.reset":         "latest",
		"client.id":                 "confluent_kafka_go_client",
		"heartbeat.interval.ms":     3000,
		"session.timeout.ms":        30000,
		"max.poll.interval.ms":      120000,
		"fetch.max.bytes":           1024000,
		"max.partition.fetch.bytes": 256000,
		"enable.auto.commit":        false, // 关闭自动提交偏移量（手动提交）
		"enable.auto.offset.store":  false, // 关闭自动存储偏移量（手动存储）
		"group.id":                  group}

	consumer, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return err
	}

	defer consumer.Close()
	flushTicker := time.NewTicker(batchInterval)
	defer flushTicker.Stop()

	// 订阅的Topics列表
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return err
	}

	logrus.Infof("Consumer inited.")

	totalReadMsg := 0

	msgChan := make(chan *kafka.Message, 1)
	defer close(msgChan)
	curSize := 0
	dataBuffer := &bytes.Buffer{}
	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				logrus.Warnf("Consumer read message error: %v (%v)\n", err, msg)
			} else {
				msgChan <- msg
				totalReadMsg++
			}

		}
	}()

	for {
		select {
		case <-flushTicker.C:
			// 时间触发器 触发batch
			if dataBuffer.Len() == 0 {
				continue
			}
			// 执行批量数据处理逻辑
			//time.Sleep(1 * time.Second)
			s := dataBuffer.String()
			//println(s)
			logrus.WithField("Trigger", "TIME").Infof("totalReadMsg: %d curSize: %d batch_size: %d \n", totalReadMsg, curSize, len(strings.Split(s, "\n")))

			// 重新攒批
			commit, err := consumer.Commit()
			logrus.WithField("Trigger", "TIME").Infof("Commit offset: %v", commit)
			if err != nil {
				logrus.WithField("Trigger", "TIME").WithField("Trigger", "TIME").Errorln("Failed to commit offset.")
				break
			}
			curSize = 0
			dataBuffer.Reset()
		case message, ok := <-msgChan:
			if !ok {
				break
			}
			dataBuffer.Write(message.Value)
			_, err := consumer.StoreMessage(message)
			if err != nil {
				logrus.Errorf("Failed to store message %v", message)
				break
			}
			curSize++
			if curSize < maxBatchSize {
				dataBuffer.Write([]byte(LineDelimiter))
			} else {
				// 执行批量数据处理逻辑
				//time.Sleep(1 * time.Second)
				s := dataBuffer.String()
				//println(s)
				logrus.WithField("Trigger", "SIZE").Infof("totalReadMsg: %d curSize: %d batch_size: %d \n", totalReadMsg, curSize, len(strings.Split(s, "\n")))

				commit, err := consumer.Commit()
				logrus.WithField("Trigger", "SIZE").Infof("Commit offset: %v", commit)
				if err != nil {
					logrus.WithField("Trigger", "SIZE").Errorln("Failed to commit offset.")
					break
				}
				curSize = 0
				dataBuffer.Reset()
			}
		}
	}
	return nil
}

func ManuallyOffsetConsume() error {
	// 手动提交Offset的消费 到达一定条数做批处理  ali写的
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           "confluent-kafka-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false", // disable auto commit
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"qjj"}, nil)

	// 用于捕获信号以优雅关闭消费者
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	batchCnt := 0
	totalCnt := 0
	const batchSize = 10000

	// <string, TopicPartition>, 用于记录每个分区的最大位点
	offsetsMap := make(map[string]kafka.TopicPartition)

	for run {
		select {
		case sig := <-sigchan:
			logrus.Infof("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := c.ReadMessage(-1)
			if err == nil {
				batchCnt++
				totalCnt++
				//fmt.Printf("[cnt:%d]Message on %s: %s\n", totalCnt, msg.TopicPartition, string(msg.Value))

				// 记录位点
				key := fmt.Sprintf("%s-%d", *msg.TopicPartition.Topic, msg.TopicPartition.Partition)
				tp := msg.TopicPartition
				if existingTp, ok := offsetsMap[key]; ok {
					if msg.TopicPartition.Offset > existingTp.Offset {
						tp.Offset = msg.TopicPartition.Offset
					} else {
						tp = existingTp
					}
				}
				offsetsMap[key] = tp

				// 批量提交
				if batchCnt >= batchSize {
					fmt.Printf("[batch_cnt:%d]\n", totalCnt)
					commitOffsets(c, offsetsMap)
					offsetsMap = make(map[string]kafka.TopicPartition)
					batchCnt = 0
				}
			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				logrus.Errorf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}

	// 程序优雅关闭前提交最后一批位点
	logrus.Infoln("Committing offsets before shutting down...")
	commitOffsets(c, offsetsMap)
	return nil
}

// commitOffsets 提交当前的偏移量
func commitOffsets(c *kafka.Consumer, offsetsMap map[string]kafka.TopicPartition) {
	var offsets []kafka.TopicPartition
	for _, tp := range offsetsMap {
		tp.Offset++ // 提交下一个位点
		offsets = append(offsets, tp)
	}

	_, err := c.CommitOffsets(offsets)
	if err != nil {
		fmt.Printf("Failed to commit offsets: %v\n", err)
	} else {
		fmt.Printf("Successfully committed offsets: %v\n", offsets)
	}
}

func appendToFile(data *bytes.Buffer, path string) error {
	s := data.String()
	println(s)
	fmt.Printf("appendToFile len: %d \n", len(strings.Split(s, "\n")))
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Failed to create file %s: %s \n", path, err)
		return err
	}
	defer file.Close()
	_, err = data.WriteTo(file)
	if err != nil {
		return err
	}
	return nil
}

func clearFile(path string) error {
	return os.Truncate(path, 0)
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
