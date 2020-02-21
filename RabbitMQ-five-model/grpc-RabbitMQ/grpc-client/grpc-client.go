package main

import (
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"grpc-RabbitMQ/proto"
	"log"
)

//微服务 客户端
func main()  {
	godotenv.Load(".env")
	conn,err := grpc.Dial("localhost:8000",grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to connect server err : %s",err)
	}
	defer conn.Close()
	msg := "charlie"
	client := proto.NewMsgClient(conn)
	resp,err := client.SendMsg(context.Background(),&proto.Request{Name:msg})
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln("receive msg:"+resp.Message)
}
