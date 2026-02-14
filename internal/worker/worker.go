package worker

import (
	"context"
	"log/slog"
)

type Config struct {
	Count     int
	QueueSize int
}

type Task func(ctx context.Context) error

type Worker struct {
	jobs chan Task
}

func New(cfg Config) *Worker {
	w := &Worker{
		jobs: make(chan Task, cfg.QueueSize),
	}

	for i := 0; i < cfg.Count; i++ {
		go w.worker()
	}
	return w
}

func (w *Worker) worker() {
	for job := range w.jobs {
		if err := job(context.Background()); err != nil {
			slog.Error("worker task", slog.String("error", err.Error()))
		}
	}
}

func (w *Worker) Push(ctx context.Context, task Task) {
	go func() {
		select {
		case w.jobs <- task:
		case <-ctx.Done():
		}
	}()
}
