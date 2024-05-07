# 限流策略
限流又称为流量控制，通常是指限制到达系统的并发请求数。

## 漏桶
无论有多少请求，请求的速率有多大，都按照固定的速率流出，对应到系统中就是按照固定的速率处理请求。
漏桶法的关键点在于漏桶始终按照固定的速率运行，但是它并不能很好的处理有大量突发请求的场景，毕竟在某些场景下我们可能需要提高系统的处理效率，而不是一味的按照固定速率处理请求。

### 代码
使用uber开源的包`go.uber.org/ratelimit`
```go
package main

import (
	"fmt"
	"time"

	"go.uber.org/ratelimit"
)

func main() {
	rl := ratelimit.New(10) // per second

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}

	// Output:
	// 0 0
	// 1 100ms
	// 2 100ms
	// 3 100ms
	// 4 100ms
	// 5 100ms
	// 6 100ms
	// 7 100ms
	// 8 100ms
	// 9 100ms
}

```
> 核心代码

每一次都确认上一次执行时间，如果没有达到时间限制则需要等待。
```go
// Take 会阻塞确保两次请求之间的时间走完
// Take 调用平均数为 time.Second/rate.
func (t *limiter) Take() time.Time {
	t.Lock()
	defer t.Unlock()

	now := t.clock.Now()

	// 如果是第一次请求就直接放行
	if t.last.IsZero() {
		t.last = now
		return t.last
	}

	// sleepFor 根据 perRequest 和上一次请求的时刻计算应该sleep的时间
	// 由于每次请求间隔的时间可能会超过perRequest, 所以这个数字可能为负数，并在多个请求之间累加
	t.sleepFor += t.perRequest - now.Sub(t.last)

	// 我们不应该让sleepFor负的太多，因为这意味着一个服务在短时间内慢了很多随后会得到更高的RPS。
	if t.sleepFor < t.maxSlack {
		t.sleepFor = t.maxSlack
	}

	// 如果 sleepFor 是正值那么就 sleep
	if t.sleepFor > 0 {
		t.clock.Sleep(t.sleepFor)
		t.last = now.Add(t.sleepFor)
		t.sleepFor = 0
	} else {
		t.last = now
	}

	return t.last
}
```
## 令牌桶
令牌桶按固定的速率往桶里放入令牌，并且只要能从桶里取出令牌就能通过，令牌桶支持突发流量的快速处理。
对于从桶里取不到令牌的场景，我们可以选择等待也可以直接拒绝并返回。
### 代码
这里使用golang官方提供的包`"golang.org/x/time/rate"`
```go
package main

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// 每秒10个令牌  桶容量5
	limiter := rate.NewLimiter(10, 5)
	cnt := 0
	last := time.Now()
	for {
		ok := limiter.Allow()
		if ok {
			cur := time.Now()
			cnt++
			fmt.Println(cnt, cur.Sub(last))
			last = cur
		} else {
			time.Sleep(time.Microsecond * 20)
			// fmt.Println("reach limit, slow down")
		}
	}
}
```
> 核心逻辑
每次请求令牌时，根据当前时间计算是否满足设定限制。同时会记录当前限制器状态。

```go
type Limiter struct {
	mu     sync.Mutex
	limit  Limit
	burst  int
	tokens float64
	// last is the last time the limiter's tokens field was updated
	last time.Time
	// lastEvent is the latest time of a rate-limited event (past or future)
	lastEvent time.Time
}

// reserveN is a helper method for AllowN, ReserveN, and WaitN.
// maxFutureReserve specifies the maximum reservation wait duration allowed.
// reserveN returns Reservation, not *Reservation, to avoid allocation in AllowN and WaitN.
func (lim *Limiter) reserveN(t time.Time, n int, maxFutureReserve time.Duration) Reservation {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	if lim.limit == Inf {
		return Reservation{
			ok:        true,
			lim:       lim,
			tokens:    n,
			timeToAct: t,
		}
	} else if lim.limit == 0 {
		var ok bool
		if lim.burst >= n {
			ok = true
			lim.burst -= n
		}
		return Reservation{
			ok:        ok,
			lim:       lim,
			tokens:    lim.burst,
			timeToAct: t,
		}
	}

	t, tokens := lim.advance(t)

	// Calculate the remaining number of tokens resulting from the request.
	tokens -= float64(n)

	// Calculate the wait duration
	var waitDuration time.Duration
	if tokens < 0 {
		waitDuration = lim.limit.durationFromTokens(-tokens)
	}

	// Decide result
	ok := n <= lim.burst && waitDuration <= maxFutureReserve

	// Prepare reservation
	r := Reservation{
		ok:    ok,
		lim:   lim,
		limit: lim.limit,
	}
	if ok {
		r.tokens = n
		r.timeToAct = t.Add(waitDuration)

		// Update state
		lim.last = t
		lim.tokens = tokens
		lim.lastEvent = r.timeToAct
	}

	return r
}
```
# 实际使用
配合Gin框架和wrk进行实战。
## 实验代码
```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	r := gin.Default()

	limiter := rate.NewLimiter(200, 200)
	r.GET("/rate", func(ctx *gin.Context) {
		ok := limiter.Allow()
		if ok {
			ctx.JSON(http.StatusOK, nil)
		} else {
			ctx.JSON(http.StatusTooManyRequests, nil)
		}
	})

	r.Run(":8080")
}
```
## 使用wrk工具进行测试
```cmd
go-wrk -c=10 -n=10000 -m="GET" "http://127.0.0.1:8080/rate"
#### 
==========================BENCHMARK==========================
URL:                            http://127.0.0.1:8080/rate

Used Connections:               10
Used Threads:                   1
Total number of calls:          10000

===========================TIMINGS===========================
Total time passed:              4.25s
Avg time per request:           4.24ms
Requests per second:            2350.38
Median time per request:        4.01ms
99th percentile time:           10.41ms
Slowest time for request:       26.00ms

=============================DATA=============================
Total response body sizes:              40000
Avg response body per request:          4.00 Byte
Transfer rate per second:               9401.53 Byte/s (0.01 MByte/s)
==========================RESPONSES==========================
20X Responses:          1048    (10.48%)
30X Responses:          0       (0.00%)
40X Responses:          8952    (89.52%)
50X Responses:          0       (0.00%)
Errors:                 0       (0.00%)
```

