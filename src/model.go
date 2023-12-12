package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Redirect struct {
	Port        int64
	Destination string
}

func (r Redirect) String() string {
	return fmt.Sprintf("%d -> %s", r.Port, r.Destination)
}

type Redirects []Redirect

// Set implements flag.Value.
func (r *Redirects) Set(str string) error {
	child, err := ParseRedirect(str)
	if err != nil {
		return err
	}
	*r = append(*r, child)
	return nil
}

// String implements flag.Value.
func (r *Redirects) String() string {
	var b strings.Builder
	for _, child := range *r {
		fmt.Fprint(&b, child)
		b.WriteString("\n")
	}
	return b.String()
}

var _ flag.Value = &Redirects{}

func ParseRedirect(str string) (Redirect, error) {
	parts := strings.SplitN(str, ":", 2)
	if len(parts) != 2 {
		fmt.Errorf("Failed to parse %s", str)
	}
	port := parts[0]
	url := parts[1]
	portNum, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		fmt.Errorf("Could not parse redirection: %s is not a number", port)
	}
	return Redirect{
		Port:        portNum,
		Destination: url,
	}, nil
}
