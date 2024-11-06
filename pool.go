package pool

import (
	"context"
	"log/slog"
)

type TJobPayload interface{}
type TJobFunc func(ctx context.Context, log *slog.Logger, payload TJobPayload) error

// Create new pool with logger and id
func New(log *slog.Logger, id string) *TPool {
	return &TPool{log: log, id: id}
}

// Pool
type TPool struct {
	id      string
	workers []*TWorker
	ch      chan TJobPayload
	log     *slog.Logger
}

// Start the pool with N workers and job handler
func (p *TPool) StartWithCancel(ctx context.Context, n int, jobfunc TJobFunc) {
	p.ch = make(chan TJobPayload)
	for i := 0; i < n; i++ {
		p.workers = append(p.workers, (&TWorker{id: i, log: p.log, pool: p.id}).new(ctx, p.ch, jobfunc))
	}
	p.log.Debug("pool is ready", "pool id", p.id, "total workers", len(p.workers))
}

// Add job to the pool
func (p *TPool) AddJob(payload interface{}) {
	p.ch <- payload
}

// Worker
type TWorker struct {
	id   int    // worker id
	pool string // pool id
	log  *slog.Logger
}

// Create new worker with context, payload channel and job handler
func (w *TWorker) new(ctx context.Context, ch chan TJobPayload, jobfunc TJobFunc) *TWorker {
	w.log = w.log.With(slog.Group("pool", slog.String("id", w.pool), slog.Int("worker", w.id)))
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.log.Debug("worker stopped")
				return
			case payload := <-ch:
				w.log.Debug("worker took the job")
				err := jobfunc(ctx, w.log, payload)
				if err != nil {
					w.log.Error("worker finished the job", "error", err)
				} else {
					w.log.Debug("worker finished the job")
				}

			}
		}
	}()
	w.log.Debug("worker started")
	return w
}
