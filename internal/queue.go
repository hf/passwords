package internal

import (
	"context"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	hashRequests chan func()
	durations    chan time.Duration

	workerThreads  int
	submittedCount uint64
	completedCount uint64

	durationMovingAverage4 int64
)

func init() {
	submittedCount = 0
	completedCount = 0
	durationMovingAverage4 = 0

	workerThreads = runtime.GOMAXPROCS(0)

	hashRequests = make(chan func())
	durations = make(chan time.Duration, workerThreads)

	for i := 0; i < workerThreads; i += 1 {
		go func() {
			for fn := range hashRequests {
				fn()
			}
		}()
	}

	go func() {
		i := 0
		values := []time.Duration{0, 0, 0, 0}

		for duration := range durations {
			values[i%len(values)] = duration
			i += 1

			ma := int64(0)

			for _, v := range values {
				ma += int64(v)
			}

			ma = ma / int64(len(values))

			atomic.StoreInt64(&durationMovingAverage4, ma)
		}
	}()
}

func SubmitRequest(ctx context.Context, threads int, fn func()) {
	if threads < 1 {
		go fn()
	} else {
		atomic.AddUint64(&submittedCount, 1)

		done := make(chan struct{})

		reserve := func() {
			<-done
		}

		work := func() {
			defer close(done)

			start := time.Now()
			fn()
			end := time.Now()

			select {
			case <-ctx.Done():
				// work was cancelled, not accounting for duration
				return

			default:
				// continue
			}

			durations <- end.Sub(start)
			atomic.AddUint64(&completedCount, 1)
		}

		reserveThreads := workerThreads
		if threads < reserveThreads {
			reserveThreads = threads
		}

		for i := 0; i < reserveThreads-1; i += 1 {
			select {
			case hashRequests <- reserve:
				// do nothing

			case <-ctx.Done():
				atomic.AddUint64(&submittedCount, ^uint64(0))
				close(done)
				return
			}
		}

		select {
		case hashRequests <- work:
			// do nothing

		case <-ctx.Done():
			atomic.AddUint64(&submittedCount, ^uint64(0))
			close(done)
			return
		}
	}
}

func NumSubmitted() uint64 {
	return atomic.LoadUint64(&submittedCount)
}

func NumCompleted() uint64 {
	return atomic.LoadUint64(&completedCount)
}

func NumOutstanding() uint64 {
	s := NumSubmitted()
	c := NumCompleted()

	if c >= s {
		return 0
	}

	return s - c
}

func NumWorkerThreads() int {
	return workerThreads
}

func DurationMovingAverage4() time.Duration {
	return time.Duration(atomic.LoadInt64(&durationMovingAverage4))
}
