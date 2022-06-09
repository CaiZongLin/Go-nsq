package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nsqio/go-nsq"
)

const (
	host     = "127.0.0.1"
	database = "vending"
	user     = "search"
	password = "123456"
	root     = "root"
	root_pwd = "s850429s"
)

func searchProduct() (string, string, string) {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	id := 5
	var product string
	var price string
	var inventory string

	err = db.QueryRow("select name,price,inventory FROM product_info where id=?", id).Scan(&product, &price, &inventory)
	return product, price, inventory
}

func main() {
	//建立Producer
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}
	a, b, c := searchProduct()
	messageBody := []byte(a)
	messageBody2 := []byte(b)
	messageBody3 := []byte(c)
	topicName := "GONSQ_TOPIC"
	//發佈到定義好的topic
	err = producer.Publish(topicName, messageBody)
	err = producer.Publish(topicName, messageBody2)
	err = producer.Publish(topicName, messageBody3)
	if err != nil {
		log.Fatal(err)
	}

	producer.Stop()
}
