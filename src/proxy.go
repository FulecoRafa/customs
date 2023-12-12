package main

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

func ListenAndLog(ctx context.Context, wg *sync.WaitGroup, r Redirect) {
    for {
        slog.Info("Listening for requests", "Redirect", r.String())
        select {
        case <-ctx.Done():
            wg.Done()
            return
        default:
        }
        time.Sleep(2 * time.Second)
    }
}
