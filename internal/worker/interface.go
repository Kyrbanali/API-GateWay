package worker

import "context"

type WorkerProvider interface {
	worker()
	Push(ctx context.Context, task Task)
}
