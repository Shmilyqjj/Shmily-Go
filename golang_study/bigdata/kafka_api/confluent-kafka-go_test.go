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
