package pool

import (
	"context"
	"log/slog"
)

type TJobPayload interface{}
type TJobFunc func(log *slog.Logger, payload TJobPayload) int

type TWorker struct {
	id     int
	poolID string
	log    *slog.Logger
}

func (w *TWorker) new(ctx context.Context, ch chan TJobPayload, jobfunc TJobFunc) *TWorker {
	w.log = w.log.With(slog.Group("pool", slog.String("id", w.poolID), slog.Int("worker", w.id)))
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.log.Debug("stopped")
				return
			case payload := <-ch:
				w.log.Debug("took the job")
				ret := jobfunc(w.log, payload)
				w.log.Debug("finished the job", "code", ret)
			}
		}
	}()
	w.log.Debug("started")
	return w
}

type TPool struct {
	id      string
	workers []*TWorker
	ch      chan TJobPayload
	log     *slog.Logger
}

// Create new pool with logger and id
func New(log *slog.Logger, id string) *TPool {
	return &TPool{log: log, id: id}
}

func (p *TPool) AddJob(payload interface{}) {
	p.ch <- payload
}

// Start the pool with N workers and job handler
func (p *TPool) StartWithCancel(ctx context.Context, n int, jobfunc TJobFunc) {
	p.ch = make(chan TJobPayload)
	for i := 0; i < n; i++ {
		p.workers = append(p.workers, (&TWorker{id: i, log: p.log, poolID: p.id}).new(ctx, p.ch, jobfunc))
	}
	p.log.Debug("pool is ready", "total workers", len(p.workers))
}
