package main

import (
	"context"
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

func MyJobHandler(log *slog.Logger, payload pool.TJobPayload) int {
	t := rand.Intn(10000)
	time.Sleep(time.Duration(t) * time.Millisecond)
	result := payload.(TPayload).N * 10
	log.Info("calculated", "result", result)
	return 0
}

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	mypool := pool.New(log, "1")
	ctx, cancel := context.WithCancel(context.Background())
	mypool.StartWithCancel(ctx, 3, MyJobHandler)
	for i := 0; i < 10; i++ {
		var payload TPayload
		payload.N = i
		mypool.AddJob(payload)
	}
	time.Sleep(10 * time.Second)
	fmt.Println("Main is stopping...")
	cancel()
	time.Sleep(10 * time.Second)
	fmt.Println("Main is stopped")
}
