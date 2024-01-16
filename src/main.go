package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/FulecoRafa/customs/lib"
)

func setLogger() {
	var handler slog.Handler
	var options *slog.HandlerOptions
	if isDebug {
		lvl := new(slog.LevelVar)
		lvl.Set(slog.LevelDebug)
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

var ports lib.Redirects

var logFormat string

var configFilePath string

func init() {
	flag.Var(&ports, "r", "Lists of ports redirecting to URLs in format 'port:url'")
	flag.Var(&ports, "redirect", "Lists of ports redirecting to URLs in format 'port:url'")
	flag.BoolVar(&isDebug, "debug", false, "Print debug logs")
	flag.StringVar(&logFormat, "l", "kv", "Log format. One of: json; kv")
	flag.StringVar(&logFormat, "logs", "kv", "Log format. One of: json; kv")
	flag.StringVar(&configFilePath, "config", "example_config.json", "Path to config file. Defaults to '~/.config/customs/config.json'")

	LoadPlugins(configFilePath)
	// a comment

	RegisterPluginFlags()
}

func main() {
	flag.Parse()

	setLogger()

	if len(ports) == 0 {
		fmt.Println("Nothing to do")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(len(ports))
	for _, r := range ports {
		go Listen(ctx, &wg, r)
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
