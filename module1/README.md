# 1.1 for 循环修改字符串数组
```go
// 通过数组下标访问修改
strs[index] = "xxx"
```


# 1.2 基于 Channel 编写一个简单的单线程生产者消费者模型
* 使用带缓冲的 channel
* 起两个 goroutine 分别生产和消费