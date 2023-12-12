package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"slices"

)

func setDebugLevel() {
    lvl := new(slog.LevelVar)
    lvl.Set(slog.LevelDebug)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: lvl,
    })

	logger := slog.New(handler)

	slog.SetDefault(logger)
}

var isDebug bool

var ports Redirects

var outputFormat string

func init() {
	flag.StringVar(&outputFormat, "o", "http", "The output format of logs. One of: curl; http")
	flag.StringVar(&outputFormat, "output", "http", "The output format of logs. One of: curl; http")
	flag.Var(&ports, "r", "Lists of ports redirecting to URLs in format 'port:url'")
	flag.Var(&ports, "redirect", "Lists of ports redirecting to URLs in format 'port:url'")
    flag.BoolVar(&isDebug, "debug", false, "Print debug logs")
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

    if isDebug {
        setDebugLevel()
    }

	slog.Debug("Starting application", "ports", ports, "outputFormat", outputFormat)

	if len(ports) == 0 {
		fmt.Println("Nothing to do")
		return
	}

	if !checkFormat() {
		panic(fmt.Sprintf("Output Format not valid: %s", outputFormat))
	}
}
