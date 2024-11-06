package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/nghtf/pool"
)

type TPayload struct {
	N int
}

func MyJobHandler(ctx context.Context, log *slog.Logger, payload pool.TJobPayload) error {
	select {
	case <-ctx.Done():
		return errors.New("job cancelled")
	case <-time.After(time.Duration(rand.Intn(10000)) * time.Millisecond):
		log.Info("result is calculated", "result", payload.(TPayload).N*10)
		return nil
	}
}

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mypool := pool.New(log, "1")
	ctx, cancel := context.WithCancel(context.Background())
	mypool.StartWithCancel(ctx, 3, MyJobHandler)
	for i := 0; i < 10; i++ {
		var payload TPayload
		payload.N = i
		mypool.AddJob(payload)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Main is stopping...")
	cancel()
	time.Sleep(10 * time.Second)
	fmt.Println("Main is stopped")
}
