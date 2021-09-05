package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errorsCount int32 = 0
	tasksCh := make(chan Task)

	wg := sync.WaitGroup{}

	wg.Add(n)
	defer func() {
		close(tasksCh)
		wg.Wait()
	}()

	for i := 0; i < n; i++ {
		go func() {
			for task := range tasksCh {
				if err := task(); err != nil {
					atomic.AddInt32(&errorsCount, 1)
				}
			}

			wg.Done()
		}()
	}

	for _, task := range tasks {
		if int(atomic.LoadInt32(&errorsCount)) >= m {
			return ErrErrorsLimitExceeded
		}

		tasksCh <- task
	}

	return nil
}
