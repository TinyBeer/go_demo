package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/panjf2000/ants/v2"
)

const (
	FiboCnt = 30 // 计算斐波那契数数量
)

func main() {
	// 初始化协程池 容量为10
	p, err := ants.NewPool(10)
	if err != nil {
		log.Fatalf("failed to initiate goroutine pool, err: %v", err)
	}
	// 运行结束后释放资源
	defer p.Release()

	cnt := FiboCnt
	var wg sync.WaitGroup // 使用WaitGroup等待所有任务完成后进行打印
	wg.Add(cnt)
	list := make([]int, cnt)
	for i := range cnt {
		// 提交计算任务
		p.Submit(func(n int) func() {
			// 由于Submit仅接受 func() 类型的任务，这里使用闭包方式提交计算参数
			return func() {
				log.Printf("开始计算第%d个斐波那契数\n", i+1) // 任务开始 用于观测任务执行过程
				list[i] = fibo(i)
				log.Printf("第%d个斐波那契数: %d\n", i+1, list[i]) // 打印计算结果 用来观测任务执行过程
				wg.Done()
			}
		}(i))
	}
	wg.Wait() // 等待任务结束

	fmt.Println(list) // 输出计算结果
}

// 使用递归方式计算斐波那契数
func fibo(n int) int {
	if n < 2 {
		return 1
	}
	return fibo(n-1) + fibo(n-2)
}
