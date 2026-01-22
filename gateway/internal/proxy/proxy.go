package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Upstream struct {
	proxy       *httputil.ReverseProxy
	stripPrefix string
	rewrite     func(*http.Request)
}

func NewUpstream(target string, stripPrefix string, rewrite func(*http.Request)) (*Upstream, error) {
	parsed, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(parsed)
	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, proxyErr error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, `{"error":"upstream_unavailable","message":"upstream unavailable"}`)
	}

	return &Upstream{
		proxy:       reverseProxy,
		stripPrefix: stripPrefix,
		rewrite:     rewrite,
	}, nil
}

func (u *Upstream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if u.stripPrefix != "" {
		trimmed := strings.TrimPrefix(r.URL.Path, u.stripPrefix)
		if trimmed == "" {
			trimmed = "/"
		}
		if !strings.HasPrefix(trimmed, "/") {
			trimmed = "/" + trimmed
		}
		r.URL.Path = trimmed
	}

	if u.rewrite != nil {
		u.rewrite(r)
	}

	u.proxy.ServeHTTP(w, r)
}
