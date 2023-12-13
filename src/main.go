package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
)

func setLogger() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
    var handler slog.Handler
    var options *slog.HandlerOptions
    if isDebug {
        options = &slog.HandlerOptions{
            Level: lvl,
        }
    }
    switch logFormat {
    case "kv":
        handler = slog.NewTextHandler(os.Stdout, options)
    case "json":
        handler = slog.NewJSONHandler(os.Stdout, options)
    default:
        panic(fmt.Sprintf("Unsupported log format: %s", logFormat))
    }

	logger := slog.New(handler)

	slog.SetDefault(logger)
}

var isDebug bool

var ports Redirects

var outputFormat string

var logFormat string

func init() {
	flag.StringVar(&outputFormat, "o", "http", "The output format of logs. One of: curl; http")
	flag.StringVar(&outputFormat, "output", "http", "The output format of logs. One of: curl; http")
	flag.Var(&ports, "r", "Lists of ports redirecting to URLs in format 'port:url'")
	flag.Var(&ports, "redirect", "Lists of ports redirecting to URLs in format 'port:url'")
	flag.BoolVar(&isDebug, "debug", false, "Print debug logs")
    flag.StringVar(&logFormat, "l", "kv", "Log format. One of: json; kv")
    flag.StringVar(&logFormat, "logs", "kv", "Log format. One of: json; kv")
}

var formats = []string{
	"curl",
	"http",
}

func checkFormat() bool {
	return slices.Contains(formats, outputFormat)
}

func main() {
	flag.Parse()

    setLogger()

	slog.Debug("Starting application", "ports", ports, "outputFormat", outputFormat)

	if len(ports) == 0 {
		fmt.Println("Nothing to do")
		return
	}

	if !checkFormat() {
		panic(fmt.Sprintf("Output Format not valid: %s", outputFormat))
	}

	ctx, cancel := context.WithCancel(context.Background())
    var wg sync.WaitGroup

    wg.Add(len(ports))
	for _, r := range ports {
		go ListenAndLog(ctx, &wg, r)
	}

	// Listen for interrupt
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer close(signalChan)
	go func() {
        <-signalChan
        cancel()
        return
	}()
    wg.Wait()
}
