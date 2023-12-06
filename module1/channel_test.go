package module1

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	/*
		基于 Channel 编写一个简单的单线程生产者消费者模型：

		队列：
		队列长度 10，队列元素类型为 int
		生产者：
		每 1 秒往队列中放入一个类型为 int 的元素，队列满时生产者可以阻塞
		消费者：
		每一秒从队列中获取一个元素并打印，队列为空时消费者阻塞
	*/

	ch := make(chan int, 10)
	defer close(ch)

	baseContext := context.Background()
	ctx, cancel := context.WithTimeout(baseContext, 5*time.Second)
	defer cancel()

	go producer(ch)
	go consumer(ch)

	<-ctx.Done()
}

func producer(ch chan<- int) {
	ticker := time.NewTicker(1 * time.Second)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		select {
		case <-ticker.C:
			n := rand.Intn(100)
			ch <- n
		}
	}
}

func consumer(ch <-chan int) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			n := <-ch
			fmt.Printf("receive: %d, now is %v\n", n, time.Now())
		}
	}
}
