package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

type myMessageHandler struct{}

func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
	log.Printf("Got message: %v", m)
	log.Printf("change to srting: %v", string(m.Body))
	return nil
}

func main() {
	//建立Consumer
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("GONSQ_TOPIC", "channel", config)
	if err != nil {
		log.Fatalf("Failed to create consumer, %v", err)
	}

	//新增handler處理收到訊息時動作
	consumer.AddHandler(&myMessageHandler{})

	//連線NSQD
	err = consumer.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Fatal(err)
	}

	//卡住，不要讓main.go執行完就結束
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()
}
