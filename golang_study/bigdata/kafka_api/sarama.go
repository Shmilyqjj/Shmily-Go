package kafka_api

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var kafkaBrokers = "broker1:9092,broker2:9092,broker3:9092"

func SaramaConsumer() {
	// sarama消费者
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer([]string{kafkaBrokers}, nil)
	if err != nil {
		fmt.Printf("Failed to start consumer: %v \n", err)
		return
	}
	partitionList, err := consumer.Partitions("test") //获得该topic所有的分区
	if err != nil {
		fmt.Println("Failed to get the list of partition:, ", err)
		return
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("test", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("Failed to start consumer for partition %d: %v \n", partition, err)
			return
		}
		wg.Add(1)
		go func(sarama.PartitionConsumer) { //为每个分区开一个go协程去取值
			for msg := range pc.Messages() { //阻塞直到有值发送过来，然后再继续等待
				fmt.Printf("Partition=%d, Offset=%d, key=%s, value=%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
			defer pc.AsyncClose()
			wg.Done()
		}(pc)
	}
	wg.Wait()
}

func SaramaSyncProducer() {
	// sarama同步生产者
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //赋值为-1：这意味着producer在follower副本确认接收到数据后才算一次发送完成。
	config.Producer.Partitioner = sarama.NewRandomPartitioner //写到随机分区中，默认设置8个分区
	config.Producer.Return.Successes = true
	msg := &sarama.ProducerMessage{}
	msg.Topic = `test`
	msg.Value = sarama.StringEncoder("Hello World!")
	client, err := sarama.NewSyncProducer([]string{kafkaBrokers}, config)
	if err != nil {
		fmt.Println("producer close err, ", err)
		return
	}
	defer client.Close()
	pid, offset, err := client.SendMessage(msg)

	if err != nil {
		fmt.Println("send message failed, ", err)
		return
	}
	fmt.Printf("分区ID:%v, offset:%v \n", pid, offset)
}

func SaramaAsyncProducer() {
	// sarama异步生产者
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true //必须有这个选项
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewAsyncProducer(strings.Split(kafkaBrokers, ","), config)
	defer p.Close()
	if err != nil {
		return
	}
	//这个部分一定要写，不然通道会被堵塞
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					logrus.Errorln(err)
				}
			case <-success:
			}
		}
	}(p)
	for {
		v := "async: " + strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000))
		fmt.Fprintln(os.Stdout, v)
		msg := &sarama.ProducerMessage{
			Topic: "test",
			Value: sarama.ByteEncoder(v),
		}
		p.Input() <- msg
		time.Sleep(time.Second * 1)
	}
}
