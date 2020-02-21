package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"grpcTest/grpc-helloworld/helloworld"
	"log"
	"strconv"
	"time"
)

func main()  {
			godotenv.Load(".env")
			conn,err := grpc.Dial("127.0.0.1:8000",grpc.WithInsecure())
			if err != nil {
				log.Println("建立连接失败")
			}
			defer conn.Close()
			client := helloworld.NewHelloWorldClient(conn)
			var name = make([]string,0)
			for i := 0; i<100; i++ {
				name = append(name, strconv.Itoa(i))
			}
				//设置客户端超时时间
				ctx,cancel := context.WithTimeout(context.Background(),time.Duration(1)*time.Minute)
				defer cancel()
				resp,err := client.SayHello(ctx,&helloworld.HelloRequest{Name:fmt.Sprintf("%s",name)})
				log.Println("客户端发送成功")
				if err != nil {
					log.Println(err)
				}
				log.Println("服务端反馈回来的信息："+resp.Msg)
}
