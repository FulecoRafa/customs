package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/FulecoRafa/customs/lib"
)

var idLock sync.Mutex
var nextId int64 = 0

func GetId() int64 {
	idLock.Lock()
	defer idLock.Unlock()
	id := nextId
	nextId++
	return id
}

func StringRequest(req *http.Request) string {
	str := "Could not dump request"
	switch outputFormat {
	case "http":
		bytes, err := httputil.DumpRequest(req, true)
		if err != nil {
			break
		}
		str = string(bytes)
	}
	return str
}

func StringResponse(resp *http.Response) string {
    str := "Could not dump response"
    switch outputFormat {
    case "http":
        bytes, err := httputil.DumpResponse(resp, true)
        if err != nil {
            break
        }
        str = string(bytes)
    }
    return str
}

func CopyHeaders(src http.Header, dst *http.Header) {
    for headingName, headingValues := range src {
        for _, value := range headingValues {
            dst.Add(headingName, value)
        }
    }
}

type proxy struct {
	redirect lib.Redirect
}

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func (p *proxy) DropHopHeaders(head *http.Header) {
	for _, header := range hopHeaders {
		head.Del(header)
	}
}

func (p *proxy) SetProxyHeader(req *http.Request) {
    headerName := "X-Forwarded-for"
    target := p.redirect.Destination
    if prior, ok := req.Header[headerName]; ok {
        // Not first proxy, append
        target = strings.Join(prior, ", ") + ", " + target
    }
    req.Header.Set(headerName, target)
}

// ServeHTTP implements http.Handler.
func (p *proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	id := GetId()
	reqStr := StringRequest(req)
	slog.Info(reqStr, "id", id)

	// Remove original URL for redirect
	req.RequestURI = ""

    // Set URL accordingly
    req.URL.Host = p.redirect.Destination
    if req.TLS == nil {
        req.URL.Scheme = "http"
    }else{
        req.URL.Scheme = "https"
    }

	// Remove connectino headers
	// (will be replaced by redirect client)
	p.DropHopHeaders(&req.Header)

	// Register Proxy Request
    p.SetProxyHeader(req)

    // Resend request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(rw, "Server Error: Redirect failed", http.StatusInternalServerError)
        slog.Error("Error redirecting request", "id", id, "error", err)
    }
    defer resp.Body.Close()

    respStr := StringResponse(resp)
    slog.Info(respStr, "id", id)

    // Once again, remove connection headers
    p.DropHopHeaders(&resp.Header)

    // Prepare and send response
    CopyHeaders(rw.Header(), &resp.Header)
    rw.WriteHeader(resp.StatusCode)
    if _, err = io.Copy(rw, resp.Body); err != nil {
        slog.Error("Error writing response", "id", id, "error", err)
    }
}

var _ http.Handler = &proxy{}

func Listen(ctx context.Context, wg *sync.WaitGroup, r lib.Redirect) {
    p := &proxy {
        redirect: r,
    }
    srvr := http.Server{
        Addr: fmt.Sprintf(":%d", r.Port),
        Handler: p,
    }
    go func() {
        defer wg.Done()
        slog.Info("Listening for requests", "Redirect", r.String())
        if err := srvr.ListenAndServe(); err != nil {
            slog.Error("Server is down", "Redirect", r.String(), "Error", err)
        }
    } ()
    defer srvr.Shutdown(context.TODO())

	<-ctx.Done()
}
