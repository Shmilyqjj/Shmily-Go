package kafka_api

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	_ "github.com/segmentio/kafka-go"
	"strings"
	"sync"
	"sync/atomic"
)

func Consumer() {
	brokers := "localhost:9092"
	topic := "qjj"
	partitionIds := []int{0, 1, 2, 3}
	readers := make(map[int]*kafka.Reader, len(partitionIds))

	for _, partitionId := range partitionIds {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   strings.Split(brokers, ","),
			Topic:     topic,
			Partition: partitionId,
		})
		readers[partitionId] = reader
	}

	var wg sync.WaitGroup
	var total int64
	total = 0
	for partitionId, partitionReader := range readers {
		// 每个分区一个reader
		partitionReader.SetOffset(699930) // 手动管理offset
		wg.Add(1)
		go func(partitionId int, reader *kafka.Reader) {
			defer wg.Done()
			currentMsgOffset := reader.Offset()
			var ml []*kafka.Message
			for {
				message, _ := reader.FetchMessage(context.Background())
				currentMsgOffset++
				atomic.AddInt64(&total, 1)
				ml = append(ml, &message)
				println(fmt.Sprintf("partitionId:%d MsgOffset:%d readerOffset:%d cnt:%d total:%d msg:%v", partitionId, message.Offset, reader.Offset(), len(ml), total, message))

				err := reader.SetOffset(currentMsgOffset)
				if err != nil {
					println(err)
				}
			}
		}(partitionId, partitionReader)
	}
	wg.Wait()

}

func ConsumerChan() {
	brokers := "localhost:9092"
	topic := "qjj"
	partitionIds := []int{0, 1, 2, 3}
	readers := make(map[int]*kafka.Reader, len(partitionIds))

	for _, partitionId := range partitionIds {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   strings.Split(brokers, ","),
			Topic:     topic,
			Partition: partitionId,
		})
		readers[partitionId] = reader
	}

	var wg sync.WaitGroup
	var total int64
	total = 0
	totalConsume := 0
	for partitionId, partitionReader := range readers {
		wg.Add(1)
		go func(partitionId int, partitionReader *kafka.Reader) {
			defer wg.Done()
			currentMsgOffset := 699930
			partitionReader.SetOffset(int64(currentMsgOffset))
			messagesChan := make(chan kafka.Message, 1)
			defer close(messagesChan)

			go func(partitionId int, reader *kafka.Reader) {
				for {
					message, _ := reader.FetchMessage(context.Background())
					messagesChan <- message

					totalConsume++
					println("totalConsume=", totalConsume)
				}
			}(partitionId, partitionReader)

			var ml []*kafka.Message
			for {
				select {
				case message := <-messagesChan:
					currentMsgOffset++
					ml = append(ml, &message)
					total++
					println(fmt.Sprintf("partitionId:%d MsgOffset:%d readerOffset:%d cnt:%d total:%d msg:%v", partitionId, message.Offset, partitionReader.Offset(), len(ml), total, message))
				}
				err := partitionReader.SetOffset(int64(currentMsgOffset))
				if err != nil {
					println(err)
				}
			}
		}(partitionId, partitionReader)

	}
	wg.Wait()
}

func GroupConsumer() {

}

func PartitionConsumer() {

}

func Producer() {

}
