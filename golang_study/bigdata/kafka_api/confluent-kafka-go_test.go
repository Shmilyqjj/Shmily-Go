package kafka_api

import "testing"

func TestDoConsume(t *testing.T) {
	err := DoConsume()
	if err != nil {
		panic(err)
	}
}

func TestDoProduce(t *testing.T) {
	DoProduce()
}

//func TestManuallyOffsetConsumer(t *testing.T) {
//	err := ManuallyOffsetConsumer()
//	if err != nil {
//		panic(err)
//	}
//}

func TestUseKafkaAdminClient(t *testing.T) {
	UseKafkaAdminClient()
}
