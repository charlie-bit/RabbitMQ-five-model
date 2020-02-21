# RabbitMQ-five-model

# Mac安装RabbitMQ
1.	在终端输入/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
2.	brew install rabbitmq
3.	vim ~/.bash_profile
4.	加入:/usr/local/sbin
5.	rabbitmq-server 启动rabbitmq服务器
对于 https://www.rabbitmq.com/getstarted.htmlRabbitMQ 官方文档介绍
首先https://www.rabbitmq.com/download.html通过官方文档下载 或者通过docker容器部署https://registry.hub.docker.com/_/rabbitmq/
github上面有rabbbitMQ的官方介绍
https://github.com/rabbitmq/rabbitmq-tutorials
# 最简单的模式
# 点对点模式：
RabbitMQ是一个消息的代理（broker），它接收并且传递（forwards）信息，可以认为它就是一个邮局。当你要寄一封邮件，你放进邮箱。你可以肯定的是快递员最终会派送给收件人（recipient）。在这个例子中，RabbitMQ就是邮局，邮箱，快递员的集成。
不同的是，Rabbitmq不会处理纸张二十接收、储存二进制的信息数据
P 代表生产者（发送信息的人）。
Queue就是邮箱，被本地记忆和磁盘（disk）限制,它是一个大型的信息缓冲容器，生产者发送信息给队列，消费者者可以从这个队列里面接收信息。
C 代表消费者（接收信息的人）
这三者可以在一台服务器上也可以不在

# Work模式 https://www.rabbitmq.com/tutorials/tutorial-two-go.html
新建一个工作队列用来分配工作任务给多个工作者这种模式主要是用来避免有密集（intensive）的任务需要等待分配，安排这个任务等下执行，我们将任务视为（encapsulate）一个信息发送给消息队列，一个工作进程将会跳出任务体（出栈），最终会执行这个任务。当运行许多工作者时，这些任务将会分配给它们。这个概念将会非常的有用在网页应用中，在短链接http请求时，可以处理复杂的任务。Preparation：发送字符串代表复制结构体，因为没有真实的任务像一些重新定义大小的图片或者一些上传的pdf文件。因此使用time.sleep函数来模仿任务繁忙，使用多个小圆点作为任务复杂度，每个小圆点将作为一个任务时间秒

# 循环调度：
使用任务队列的优点是我们可以并行（parallelise）工作,可以增加更多的工作者来缓解处理任务的压力。默认情况下，RabbitMQ将每个信息发送给下一个工作者，每个工作者将获得相同数量的信息。根据现在的代码情况，一旦生产者将信息发送给消费者，这个任务数据就会自动被杀死，也就是数据会丢失。
为了保证消息确认，消费者死亡时，将重新传递消息。
队列持久化操作
msgs，err：= ch.Consume（   
q.Name，//队列   
“”，      //消费者   
false，   //自动确认   
false，   //排他   
false，   // no-local    
false，   // no-wait   
 nil，     // args  ）
err = ch.Publish(   
"",           // exchange   
q.Name,       // routing key   
false,        // mandatory   
false,   
amqp.Publishing {     
DeliveryMode: amqp.Persistent,     
ContentType:  "text/plain",     
Body:         []byte(body), })
# 公平分配
消息如果是几秒就完成时，一个工作者保持忙碌状态，另外一名工作者几乎什么都不做，这些以上两种模式都没有控制。因此我们需要设置，一名工作处理任务的量。在客户端消费者前面设置也就是channel.Consume
err = ch.Qos(                 
1,     // prefetch count                 
0,     // prefetch size                 
false, // global         )
                                                                                                                                        
# 发布-订阅模式 https://www.rabbitmq.com/tutorials/tutorial-three-go.html
在上面的集中模式中，我们假设通过一个工作低劣，每个任务恰恰分配给一个工作者。在这个部分，我们将做一些不同的事情，这一部分就是发布-订阅。
为了说明这一部分我们将建一个简单的日志系统，分为两部分，第一部分发送日志信息，第二部分接收并且打印他们。
在日志系统中，每次运行接收项目时都会得到信息，通过这种方法，我们将运行接收者项目将日志写到磁盘里面。同时，运行另外一个接受者项目在控制台界面看日志。最基本的，发布日志信息将会传播给所有的接收者。
核心思想：在RabbitMQ信息模型中，生产者从不发送任何消息给队列。生产者只将信息发送到交换机中，交换机是个非常简单的事情。一边接收者从生产者那边接收信息，另外一边，生产者将信息推送给队列。交换机必须准确的知道怎么样解决接受的信息。
交换机有四种类型：direct、topic、headers、fanout.
Fanout交换机类型：
err = ch.ExchangeDeclare(   
"logs",   // name   
"fanout", // type   
true,     // durable  
 false,    // auto-deleted   
false,    // internal   
false,    // no-wait   
nil,      // arguments )
默认交换器：
err = ch.Publish(   
"",     // exchange   
q.Name, // routing key   
false,  // mandatory   
false,  // immediate   
amqp.Publishing{    
 ContentType: "text/plain",     
Body:        []byte(body), })
将日志交换机替换默认交换器
err = ch.ExchangeDeclare(   
"logs",   // name   
"fanout", // type   
true,     // durable   
false,    // auto-deleted   
false,    // internal   
false,    // no-wait   
nil,      // arguments ) 
failOnError(err, "Failed to declare an exchange")  body := bodyFrom(os.Args) 
err = ch.Publish(   
"logs", // exchange   
"",     // routing key   
false,  // mandatory   
false,  // immediate   
amqp.Publishing{           
ContentType: "text/plain",           
Body:        []byte(body),   })
暂时性队列：
之前两个版本，使用了特定的队列，生产者发布了特定的名字，接收者接收信息需要consume特定的队列。
当我们连接RabiitMQ，需要使用一个新的，空的队列，否则会报错，应为本地Rabbitmq服务器已经产生了一种类型的队列，必须定义一个新的队列，还可以重启RabbitMQ服务器。
定义一个交换器和一个空的队列，需要绑定队列的名字，交换器是唯一的。
err = ch.QueueBind(   
q.Name, // queue name   
"",     // routing key   
"logs", // exchange   
false,   
nil, )
channle.QueueDeclare方法返回的时候，队列实例化包含的是一个随机的队列，它的名字是被RabbitMQ默认的。当连接关闭的时候，队列将会被删除，因为它声明的时候就是独立的。
我们创建了一个fanout类型的交换器，需要告诉交换器，我们已经将信息发送给了队列，他们之间的关系被称之为绑定binding
在终端界面，输入如下指令可以或者已经存在的队列。
rabbitmqctl list_bindings

结果发现队列会同时接受相等数量的数据，在终端发现，mq会产生两个独立的队列出来。
路由模式  https://www.rabbitmq.com/tutorials/tutorial-four-go.html
在交换器与队列绑定的过程中增加一种routing-key的参数
direct交换器：
我们之前的日志系统传播所有的信息给消费者，我们在他们的复杂度上面拓展去过滤信息。例如我们希望脚本只接受严重错误的信息并写进磁盘中同时又不会浪费磁盘空间也不会有错误的信息。
我们使用fanout交换器，并没有很大的灵活性，因为他是非自主性传播，因此我们使用direct交换器。
Direct交换器路由算法（algorithm）很简单，消息只会传送到已经绑定routing-key的队列中。
direct交换器有两个绑定的队列，Q1绑定的routing-key是orange，
Q2绑定的routing-key是black和green。
消息发布给交换器，routing-key是orange的将会发送给Q1，routing-key是black、green的会发送给Q2。其它信息将会被丢弃掉。
多个绑定：
一个绑定routing-key绑定多个队列是合法的。
direct交换传播信息的方式与fanout交换器传播信息的方式是一致的。当信息由direct交换器传播给routing-key是black的队列，也就是图中的Q1、Q2。
发送日志：
将复杂的日志作为一个routing-key，这种接受脚本的方式将会得到接受信息的复杂度。
 
# topics模式：
direct模式在fanout模式上改进了程序，使得每条信息都有对应的队列进行传输，但是不能基于多重条件的路由，只是使用特定的routing-key。
在日志系统中，我们想要订阅的不仅仅基于复杂度的日志，也针对于日志的来源。
Routing-key的灵活性体现在我们不仅仅监听到了warning信息，我们还需要监听到信息的来源，为了改进我们的日志系统，我们需要使用topic模式。
Topic交换器：
信息集群发送到topic交换器里面，并没有具体的（arbitrary）routing-key，它必须是一串字符串，被小圆点所分割，这个字符串可以是任意单词，但是需要最大255字节。
绑定的路由键必须具有相同的格式，topic交换器背后的思想是和direct交换器相似的。一条信息发送到绑定routing-key队列，在topic交换器中,绑定routing-key有两条规则：
	* 可以代替一个单词
	# 可以代替0个或多个单词。


 举个例子：发送一些描述动物的数据，这些信息都会被发送到特定的routing-key（包含两个小圆点），第一个单词描述速度，第二个描述颜色，第三个描述种类。
现在有两个队列：
Q1 绑定的routing-key是：*.orange.* 
Q2 绑定的routing-key是：*.*.rabbit和lazy.#
一条信息发送给routing-key是quick.orange.rabbit将会分发给两个队列。
在日志系统中使用topic交换器，开始的时候假设routing-key是facility.severity
