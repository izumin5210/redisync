package redisync_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/izumin5210/redisync"
)

func ExampleMutex() {
	defer cleanupTestRedis()

	var wg sync.WaitGroup
	ctx := context.Background()

	mu := redisync.NewMutex(pool, "key")

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock(ctx)
			defer mu.Unlock(ctx)

			fmt.Println("start")
			time.Sleep(10 * time.Millisecond)
			fmt.Println("stop")
		}()
	}

	wg.Wait()

	// Output:
	// start
	// stop
	// start
	// stop
	// start
	// stop
	// start
	// stop
	// start
	// stop
}

func ExampleMonitor() {
	defer cleanupTestRedis()

	var wg sync.WaitGroup
	ctx := context.Background()

	monitor := redisync.NewMonitor(pool)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			monitor.Synchronize(ctx, "key", func(context.Context) error {
				fmt.Println("start")
				time.Sleep(10 * time.Millisecond)
				fmt.Println("stop")
				return nil
			})
		}()
	}

	wg.Wait()

	// Output:
	// start
	// stop
	// start
	// stop
	// start
	// stop
	// start
	// stop
	// start
	// stop
}

func ExampleOnce() {
	defer cleanupTestRedis()

	var wg sync.WaitGroup
	ctx := context.Background()

	once := redisync.NewOnce(pool)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(ctx, "key", func(context.Context) error {
				fmt.Println("Only once")
				return nil
			})
		}()
	}

	wg.Wait()

	// Output:
	// Only once
}

func ExampleScoreFilter() {
	defer cleanupTestRedis()

	ctx := context.Background()

	monitor := redisync.NewMonitor(pool)
	filter := redisync.NewScoreFilter(pool, monitor)

	for i := 10; i > 0; i-- {
		i := i

		ok, err := filter.Filter(ctx, "key", i)
		if err != nil {
			// ...
		}
		if ok {
			fmt.Println("Only once")
		}
	}

	// Output:
	// Only once
}
