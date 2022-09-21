package http

import (
	"errors"
	"strings"
)

type Request struct {
    Path string
    Method string
    Version string
    Headers map[string]string
}

// Parse parses the raw output from the TCP connection as an HTTP
// request object.
func Parse(lines []string) (*Request, error) {
    head := lines[0]
    headers := lines[1:]

    if !strings.HasSuffix(head, "HTTP/1.1") {
        return nil, errors.New("unsupported HTTP version")
    }

    req := Request{
        Headers: map[string]string{},
    }

    headParts := strings.Split(head, " ")
    req.Method = headParts[0]
    req.Path = headParts[1]
    req.Version = headParts[2]

    for _, v := range headers { 
        kv := strings.Split(v, ": ")
        req.Headers[kv[0]] = kv[1]
    }

    return &req, nil
}

