package main

import (
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"grpc-RabbitMQ/RabbitMq"
	"grpc-RabbitMQ/proto"
	"log"
	"net"
	"os"
)
type msg struct{}

func (m msg) SendMsg(ctx context.Context,in *proto.Request) (resp *proto.Response,err error) {
	msg := "hello" + in.Name
	//将回馈信息放进消息队列里面
	err = RabbitMq.PushRMQ(msg)
	if err != nil {
		log.Fatalf("send message into RMQ err:%s",err)
	}
	resp.Message = msg
	return resp,nil
}

//微服务 服务端
func main()  {
	//初始化消息队列
	err := RabbitMq.InitRMQ()
	if err != nil {
		log.Fatalf("connect rabbitmq err:%s",err)
	}
	godotenv.Load(".env")
	address := ":" + os.Getenv("server_port")
	listener,err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("fail to listen %s",err)
	}
	server := grpc.NewServer()
	proto.RegisterMsgServer(server,&msg{})
	log.Fatalln("server run successfully...")
	if err := server.Serve(listener) ; err != nil {
		log.Fatalf("fail to listen err : %s",err)
	}

}
