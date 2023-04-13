## Y-MQ

Y-MQ利用高效可靠的消息传递机制进行**异步的数据传输**，并基于数据通信进行分布式系统的集成，通过提供消息队列模型和消息传递机制，可在分布式环境下扩展进程间的通信。



- #### version 1.0（单机版本）

  实现消息队列的基本功能，采用服务发布与订阅模式。服务器发送消息确保至少发送一次，并且在消息传输中通过三个管道：

  ##### Channel设计

  宏观Channel由三个缓冲层组成：incomingChan -> msgChan -> clientChan，用户从clientChan中读取数据。

  Router（处理incomingMsg -> msgChan）和它的子协程MessagePump（处理msgChan -> clientChan）和RequeueRouter负责监听管道，管道触发（往管道里面写）则由对应的方法完成。

  （其中closeChan主要起到context包的作用，负责优雅地结束子协程）

  发送消息给服务端到达msgChan后会被添加到infilghtMessageMap中并开启一个超时处理协程Timer，当客户端从clientChan中获取到消息时，会向服务器返回确认消息并将消息从map中delete掉，并关闭Timer协程；

  如果数据长时间停留在clientChan中不被读走，则RequeueRouter中的子协程会触发Timer超时，将消息从map中delete掉，然后重新发送到incomingChan中。

  ##### 数据持久化

  在topic和channel中都做了磁盘持久化。

  msgChan是个有缓冲的管道，用来暂存消息。当发送消息并发量过大导致msgChan的缓冲区填满时，则会通过后台队列的形式存放到内存。

![](https://github.com/yanyanran/pictures/blob/main/MQversion1.png?raw=true)

- #### version 1.1（单机版本实现多订阅）

  ##### 设计预期结果

  在Channel的设计上选择将clientChan脱离开channel，分离出来的这类管道由protocol进行管理。

  一个客户端对应一个clientChan（后续还可以添加其他管道），服务端将消息发送给channel后，消息从msgChan中将消息通过负载均衡传递给clientChan中去，客户端读取消息只需要针对对应的clientChan读取即可。

  关于指定多订阅模式的实现，为了让客户端可以自由选择想要订阅的channel，可以开辟一个map存放客户端想要订阅的channel。

  客户端获取消息方式：指定管道或获取/随机获取管道消息

  ![](https://github.com/yanyanran/pictures/blob/main/newMQmodel.png?raw=true)
