package main

import "Shmily-Go/golang_study/bigdata/kafka"

func main() {
	//kafka.Consumer()
	//kafka.SyncProducer()
	go kafka.AsyncProducer()
	select {}
}
