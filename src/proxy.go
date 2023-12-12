package main

import (
	"context"
	"log/slog"
	"time"
)

func ListenAndLog(ctx context.Context, exitChan chan<- struct{}, r Redirect) {
    defer func() {exitChan <- struct{}{}}()
    for {
        slog.Info("Listening for requests", "Redirect", r.String())
        select {
        case <-ctx.Done():
            exitChan<- struct{}{}
            return
        default:
        }
        time.Sleep(2 * time.Second)
    }
}
