package metrics

import (
	"time"

	"github.com/hf/passwords/internal"
)

// NumOutstanding returns the number of outstanding (queued) password runs
// requests.
func NumOutstanding() uint64 {
	return internal.NumOutstanding()
}

// NumCompleted returns the number of successfully completed password runs
// requests.
func NumCompleted() uint64 {
	return internal.NumCompleted()
}

// NumSubmitted returns the number of submitted password hashing runs. A
// cancelled attempt decrements this number.
func NumSubmitted() uint64 {
	return internal.NumSubmitted()
}

// DurationMovingAverage4 returns the 4-point moving average of durations
// observed on completed password hashing runs.
func DurationMovingAverage4() time.Duration {
	return internal.DurationMovingAverage4()
}

// DurationQueue gives an estimate about how long it will take to clear the
// NumOutstanding password hashing runs. You can use this to preemptively
// reject new attempts if the queue is longer than some acceptable value
// (multiple seconds).
func DurationQueue() time.Duration {
	return time.Duration(internal.DurationMovingAverage4() * time.Duration(internal.NumOutstanding()))
}
