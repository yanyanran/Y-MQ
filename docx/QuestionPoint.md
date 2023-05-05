### 现有问题

1、消息有序化（预期结合卡夫卡修改[没有msgChan消息直接存磁盘]（看论文））

2、消息的延迟发送在高并发场景下，进入磁盘后延迟丢失：

​	（1）消息进入Topic，发现msgChan满了，写入Topic diskpq，此时延迟丢失；

​	（2）进入Topic发现msgChan没满，通过msgChan广发给Channel。发现Channel的msgChan满了就往Channel diskpq写，此时也会造成延迟丢失。

3、 msgChan中消息可能会丢失（没有保存、断电宕机就没了）

4、消息类型太少（延迟的还拉）