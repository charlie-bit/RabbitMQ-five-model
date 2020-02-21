package RabbitMq

import (
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

var conn *amqp.Connection
var channel *amqp.Channel
//错误输出
func Fail(msg string,err error)  {
	if err != nil {
		log.Fatalf("%s,%s",msg,err)
	}
}
//初始化消息队列服务端
func InitRMQ(){
	if channel == nil {
		var err error
		godotenv.Load(".env")
		url := "amqp://" + os.Getenv("username") + ":" +
			os.Getenv("password") + "@" + os.Getenv("rabbitmq_ip") + ":" + os.Getenv("rabbitmq_port") + "/"
		conn,err = amqp.Dial(url)
		Fail("fail to connect",err)
		channel,err = conn.Channel()
		Fail("fail to create channel",err)
	}
}
//将数据放进消息队列
func PushRMQ(queueName,msg string) error {
	queue,err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
		)
	Fail("fail to create queue",err)
	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			Headers:         nil,
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       "",
			Timestamp:       time.Time{},
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            []byte(msg),
		})
	Fail("fail to publish a message",err)
	log.Printf("successfully publish message: %s",msg)
	return nil
}

//消费消息队列
func GetRMQ(queueName string,Function func(json string)(string)) (result string,err error) {
	receiveConn,err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Printf("%s,%s","fail to connect RabbitMQ",err)
		return "",err
	}
	defer receiveConn.Close()
	channel,err := receiveConn.Channel()
	if err != nil {
		log.Printf("%s,%s","fail to open a channel",err)
		return "",err
	}
	defer channel.Close()
	queue,err := channel.QueueDeclare(queueName,true,false,false,false,nil)
	if err != nil {
		log.Printf("%s,%s","fail to declare a queue",err)
		return "",err
	}
	msg,err := channel.Consume(queue.Name,"",false,false,false,false,nil)
	if err != nil {
		log.Printf("%s,%s","fail to regigter",err)
		return "",err
	}
	listen := make(chan bool)
	go func() {
		for d := range msg {
			log.Printf("Received message is %s",d.Body)
			//取得存到消息队列里面数据 --> 应用函数
			result = Function(string(d.Body))
			d.Ack(false)
			log.Printf("successfully get result: %s",result)
			close(listen)
			return
		}
	}()
	<-listen
	return result,nil
}
