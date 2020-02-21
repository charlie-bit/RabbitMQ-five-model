package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main()  {
	conn,err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("fail to connect")
	}
	defer conn.Close()

	channel,err := conn.Channel()
	if err != nil{
		log.Println("fail to declare channel")
	}
	defer channel.Close()

	err = channel.ExchangeDeclare("logs","fanout",true,false,false,false,nil)
	if err != nil{
		log.Println("fail to declare exchange")
	}

	queue,err := channel.QueueDeclare("",false,false,true,false,nil)
	if err != nil {
		log.Println("fail to declare queue")
	}

	err = channel.QueueBind(queue.Name,"","logs",false,nil)
	if err != nil {
		log.Println("fail to bind a queue")
	}
	msgs,err := channel.Consume(queue.Name,"",false,false,false,false,nil)
	if err != nil {
		log.Println("fail to register consumer")
	}

	lis := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("receive message : %s",d.Body)
		}
	}()

	<-lis
}
