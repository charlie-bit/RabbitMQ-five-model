package main

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main()  {
	conn,err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("fail to connect rabbitmq")
	}
	defer conn.Close()

	channel,err := conn.Channel()
	if err != nil {
		log.Println("fail to open channel")
	}

	err = channel.ExchangeDeclare("logs","fanout",true,false,false,false,nil)
	if err != nil {
		log.Println("fail to declare exchange")
	}

	body := "hello..."

	err = channel.Publish("logs","",false,false,amqp.Publishing{
		Headers:         nil,
		ContentType:     "text/plain",
		ContentEncoding: "",
		DeliveryMode:    0,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Time{},
		Type:            "",
		UserId:          "",
		AppId:           "",
		Body:            []byte(body),
	})
	if err != nil {
		log.Println("fail to publish")
	}

	log.Printf("send %s successfully",body)
}
