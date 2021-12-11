package cache

import (
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestRandSource 测试线程安全和非线程安全的 rand.Source 实现。
// 因为是并发测试，所有不一定每次都成功需要多跑几次。
// 每次的结果 safe 一定要是 pass，unsafe 只要有一次是 pass 就算是测试通过。
// note: 测试结果是有缓存的,记得清除再跑或者修改并发参数,不然会用缓存。
func TestRandSource(t *testing.T) {
	const goroutines int = 10
	const loopTimes int = 1e5

	seed := time.Now().UnixNano()
	tests := []struct {
		name      string
		source    rand.Source
		wantPanic bool
	}{
		{"safe", newConcurrentRandSource(seed), false},
		{"unsafe", rand.NewSource(seed), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := func() (e error) {
				r := rand.New(tt.source)
				var panicCount int32 = (int32)(0)
				panicCountPtr := &panicCount
				wg := &sync.WaitGroup{}

				for i := 0; i < goroutines; i++ {
					go func() {
						defer func() {
							err := recover()
							if err != nil {
								atomic.AddInt32(panicCountPtr, 1)
							}
						}()

						wg.Add(1)
						defer wg.Done()

						for i := 0; i < loopTimes; i++ {
							r.Int63()
						}
					}()
				}

				wg.Wait()

				if *panicCountPtr == 0 {
					return nil
				}
				return errors.New("panic has occurred")
			}()

			if (err != nil) != tt.wantPanic {
				t.Errorf(tt.name)
			}
		})
	}
}
