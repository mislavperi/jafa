package utils

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
)

type Task interface {
	Start(ctx context.Context)
}

// Run starts all tasks and blocks until they finish. The shared context is
// cancelled when any task returns or the process receives SIGINT/SIGTERM, so
// tasks can shut down gracefully.
func Run(tasks ...Task) {
	wg := sync.WaitGroup{}
	wg.Add(len(tasks))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	for _, task := range tasks {
		go func(task Task) {
			defer wg.Done()
			defer cancel()

			task.Start(ctx)
		}(task)
	}
	wg.Wait()
	cancel()
}
