package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/FulecoRafa/customs/lib"
)

func ParseHeader(header string) (string, string, error) {
	parts := strings.SplitN(header, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		slog.Debug("Failed to parse header", "length", len(parts), "header", parts[0], "value", parts[1])
		return "", "", fmt.Errorf("Failed to parse header")
	}
	return parts[0], parts[1], nil
}

type InjectHeader struct {
	headers map[string][]string
}

// String implements flag.Value.
func (ij *InjectHeader) String() string {
	var sb strings.Builder
	for k, v := range ij.headers {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	return sb.String()
}

// Set implements flag.Value.
func (ij *InjectHeader) Set(header string) error {
	header, value, err := ParseHeader(header)
	if err != nil {
		return err
	}
	ij.headers[header] = append(ij.headers[header], value)
	return nil
}

// AddConfigFlag implements lib.CustomsPlugin.
func (ih InjectHeader) AddConfigFlag() {
	flag.Var(&ih, "header", "Headers to inject in format `header:value`")
	flag.Var(&ih, "H", "Headers to inject in format `header:value`")
}

// PostRequestHook implements lib.CustomsPlugin.
func (InjectHeader) PostRequestHook(res *http.Response, r lib.Redirect) {
}

// PreRequestHook implements lib.CustomsPlugin.
func (ij InjectHeader) PreRequestHook(req *http.Request, r lib.Redirect) {
	for header, values := range ij.headers {
		strValues := strings.Join(values, ",")
		req.Header.Set(header, strValues)
	}
}

var _ lib.CustomsPlugin = (*InjectHeader)(nil)
var _ flag.Value = (*InjectHeader)(nil)

var Plugin = InjectHeader{}

func init() {
	Plugin.headers = make(map[string][]string)
}

func main() { /* Does nothing because it's a plugin */ }
