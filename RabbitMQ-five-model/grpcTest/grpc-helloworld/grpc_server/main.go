package main

import (
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"grpcTest/RabbitMq"
	"grpcTest/controller"
	pb "grpcTest/grpc-helloworld/helloworld"
	"log"
	"net"
)

type server struct {}

func (s *server)SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply,error)  {
	msg := in.Name
	//将请求信息放进独立消息队列里面
	err := RabbitMq.PushRMQ("unique",msg)
	if err != nil {
		log.Printf("send message into RMQ err:%s",err)
		return &pb.HelloReply{Msg:"send message into RMQ err"},err
	}
	//消费队列里面的数据
	result,err := RabbitMq.GetRMQ("unique",controller.SayHello)
	if err != nil {
		return &pb.HelloReply{Msg:"get message from RMQ err"},err
	}
	return &pb.HelloReply{Msg:"successfully send message: "+result},nil
}

func main()  {
	RabbitMq.InitRMQ()
	godotenv.Load(".env")
	lis,err := net.Listen("tcp",":8000")
	if err != nil {
		log.Printf("fail to listen :%v",err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloWorldServer(s, &server{})

	log.Println("server run successfully...")

	if err := s.Serve(lis) ;err !=nil{
		log.Println(err)
	}
}
