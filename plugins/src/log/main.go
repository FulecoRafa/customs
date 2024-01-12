package main

import (
	"flag"
	"log/slog"
	"net/http"
	"net/http/httputil"

	"github.com/FulecoRafa/customs/lib"
)

type LogPlugin struct{
    outputFormat string
}

// AddConfigFlag implements lib.CustomsPlugin.
func (lp *LogPlugin) AddConfigFlag() {
	flag.StringVar(&lp.outputFormat, "o", "http", "The output format of logs. One of: curl; http")
	flag.StringVar(&lp.outputFormat, "output", "http", "The output format of logs. One of: curl; http")
    slog.Debug("Loaded flags", "outputFormat", lp.outputFormat)
}

// PostRequestHook implements lib.CustomsPlugin.
func (lp *LogPlugin) PostRequestHook(res *http.Response, r lib.Redirect) {
    str := lp.StringResponse(res)
    slog.Debug(str)
}

// PreRequestHook implements lib.CustomsPlugin.
func (lp *LogPlugin) PreRequestHook(req *http.Request, r lib.Redirect) {
    str := lp.StringRequest(req)
    slog.Debug(str)
}

var _ lib.CustomsPlugin = &LogPlugin{}

func (lp LogPlugin) StringRequest(req *http.Request) string {
	str := "Could not dump request"
	switch lp.outputFormat {
    default:
		bytes, err := httputil.DumpRequest(req, true)
		if err != nil {
			break
		}
		str = string(bytes)
	}
	return str
}

func (lp LogPlugin) StringResponse(resp *http.Response) string {
	str := "Could not dump response"
	switch lp.outputFormat {
    default:
		bytes, err := httputil.DumpResponse(resp, true)
		if err != nil {
			break
		}
		str = string(bytes)
	}
	return str
}

var Plugin = LogPlugin{}

func main() { /* Does nothing since is a plugin*/ }
