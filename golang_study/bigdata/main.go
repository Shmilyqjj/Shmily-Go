package main

import "Shmily-Go/golang_study/bigdata/kafka_api"

func main() {
	err := kafka_api.ManuallyOffsetConsumer()
	if err != nil {
		println(err)
	}

	select {}
}
