package concurrent

import (
	"sync"
)

const (
	defaultMaxConcurrency = 100
)

func Exec[T, R any](tasks []T, concurrency int, fn func(T) (R, error)) ([]R, error) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))

	if concurrency <= 0 {
		concurrency = len(tasks)
	}
	sem := make(chan struct{}, concurrency)
	results := make([]R, len(tasks))
	errs := make([]error, len(tasks))

	for i, task := range tasks {
		sem <- struct{}{}
		go func(i int, task T) {
			defer func() {
				<-sem
				wg.Done()
			}()

			result, err := fn(task)
			if err != nil {
				errs[i] = err
				return
			}
			results[i] = result
		}(i, task)
	}

	wg.Wait()

	var err error
	for i := range errs {
		if errs[i] != nil {
			err = errs[i]
			break
		}
	}

	return results, err
}

func ExecChan[T, R any](tasks <-chan T, concurrency int, fn func(T) (R, error)) ([]R, error) {
	var wg sync.WaitGroup

	if concurrency <= 0 {
		concurrency = defaultMaxConcurrency
	}
	sem := make(chan struct{}, concurrency)
	results := make([]R, 0)
	errs := make([]error, 0)

	for task := range tasks {
		sem <- struct{}{}
		wg.Add(1)
		go func(task T) {
			defer func() {
				<-sem
				wg.Done()
			}()

			result, err := fn(task)
			if err != nil {
				errs = append(errs, err)
				return
			}
			results = append(results, result)
		}(task)
	}

	wg.Wait()

	var err error
	for i := range errs {
		if errs[i] != nil {
			err = errs[i]
			break
		}
	}

	return results, err
}
