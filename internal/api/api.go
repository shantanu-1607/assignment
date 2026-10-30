package api

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/shantanu-1607/assignment/internal/service"
)

type Options struct {
	Timeout       time.Duration
	Rate          int
	MaxConcurrent int
	Logger        *slog.Logger
}

type API struct {
	estimator service.Estimator
	limiter   *limiter
	slots     chan struct{}
	timeout   time.Duration
	log       *slog.Logger
	ready     atomic.Bool
	catalog   []byte
	etag      string
}

type problem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type envelope struct {
	Data      any      `json:"data,omitempty"`
	Error     *problem `json:"error,omitempty"`
	RequestID string   `json:"request_id"`
}

func New(estimator service.Estimator, opts Options) *API {
	if opts.Timeout <= 0 {
		opts.Timeout = 2 * time.Second
	}
	if opts.Rate <= 0 {
		opts.Rate = 60
	}
	if opts.MaxConcurrent <= 0 {
		opts.MaxConcurrent = 32
	}
	if opts.Logger == nil {
		opts.Logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	catalog, _ := json.Marshal(struct {
		Data []service.Offering `json:"data"`
	}{service.Catalog()})
	sum := sha256.Sum256(catalog)
	a := &API{estimator: estimator, limiter: &limiter{entries: make(map[string]bucket), limit: opts.Rate, maxKeys: 1024, now: time.Now}, slots: make(chan struct{}, opts.MaxConcurrent), timeout: opts.Timeout, log: opts.Logger, catalog: catalog, etag: fmt.Sprintf(`"%x"`, sum)}
	a.ready.Store(true)
	return a
}

func (a *API) SetReady(ready bool) { a.ready.Store(ready) }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, status int, id, code, message string) {
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, status, envelope{Error: &problem{code, message}, RequestID: id})
}

// ServeHTTP sets response policy outside the timeout handler, preserving it on 503s.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		fail(w, 503, "unavailable", "unavailable", "Request ID generation unavailable")
		return
	}
	id := hex.EncodeToString(random[:])
	h := w.Header()
	h.Set("Content-Type", "application/json")
	h.Set("X-Request-ID", id)
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
	h.Set("Referrer-Policy", "no-referrer")
	h.Set("X-Frame-Options", "DENY")
	h.Set("Cache-Control", "no-store")
	defer func() {
		if recover() != nil {
			a.log.Error("handler panic", "request_id", id)
			fail(w, 500, id, "internal_error", "Internal server error")
		}
		a.log.Info("request completed", "request_id", id, "duration_ms", time.Since(start).Milliseconds())
	}()
	// Only GET health/readiness bypass quota/admission, so malformed health methods
	// cannot be used to enter domain logic. Health handlers do no external work.
	if (r.URL.Path == "/healthz" || r.URL.Path == "/readyz") && r.Method == http.MethodGet {
		if r.URL.Path == "/readyz" && !a.ready.Load() {
			fail(w, 503, id, "not_ready", "Server is draining")
			return
		}
		writeJSON(w, 200, envelope{Data: map[string]string{"status": "ok"}, RequestID: id})
		return
	}
	if !a.ready.Load() {
		fail(w, 503, id, "not_ready", "Server is draining")
		return
	}
	peer, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		peer = r.RemoteAddr
	}
	if !a.limiter.allow(peer) {
		h.Set("Retry-After", "60")
		fail(w, 429, id, "rate_limited", "Request quota exceeded")
		return
	}
	select {
	case a.slots <- struct{}{}:
	default:
		fail(w, 503, id, "busy", "Server at capacity")
		return
	}
	// Release inside the worker: a timeout does not imply the underlying handler exited.
	worker := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { <-a.slots }()
		defer func() {
			if recover() != nil {
				a.log.Error("worker panic", "request_id", id)
				fail(w, 500, id, "internal_error", "Internal server error")
			}
		}()
		a.route(w, r, id)
	})
	body, _ := json.Marshal(envelope{Error: &problem{"timeout", "Request deadline exceeded"}, RequestID: id})
	// TimeoutHandler bounds buffered response time and cancels the worker context.
	// Domain code must cooperate with cancellation; Go cannot forcibly stop a goroutine.
	http.TimeoutHandler(worker, a.timeout, string(body)).ServeHTTP(w, r)
}

func (a *API) route(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	switch r.URL.Path {
	case "/v1/services":
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			fail(w, 405, id, "method_not_allowed", "Use GET or HEAD")
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=60, must-revalidate")
		w.Header().Set("ETag", a.etag)
		for _, candidate := range strings.Split(r.Header.Get("If-None-Match"), ",") {
			candidate = strings.TrimSpace(candidate)
			if candidate == "*" || strings.TrimPrefix(candidate, "W/") == a.etag {
				w.WriteHeader(304)
				return
			}
		}
		w.WriteHeader(200)
		if r.Method != http.MethodHead {
			_, _ = w.Write(a.catalog)
		}
	case "/v1/estimates":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			fail(w, 405, id, "method_not_allowed", "Use POST")
			return
		}
		mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mt != "application/json" {
			fail(w, 415, id, "unsupported_media_type", "Use application/json")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		defer r.Body.Close()
		in, err := decodeInput(r.Body)
		if err != nil {
			var large *http.MaxBytesError
			if errors.As(err, &large) {
				fail(w, 413, id, "body_too_large", "Body exceeds 4096 bytes")
			} else {
				fail(w, 400, id, "invalid_json", "Expected one object with service_id and integrations; no unknown or duplicate fields")
			}
			return
		}
		out, err := a.estimator.Estimate(r.Context(), in)
		switch {
		case errors.Is(err, service.ErrInvalid):
			fail(w, 422, id, "validation_failed", "service_id must be a catalog ID; integrations must be 0..20")
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			fail(w, 503, id, "timeout", "Request deadline exceeded")
		case err != nil:
			fail(w, 500, id, "internal_error", "Internal server error")
		default:
			writeJSON(w, 200, envelope{Data: out, RequestID: id})
		}
	default:
		fail(w, 404, id, "not_found", "Route not found")
	}
}

// Token decoding rejects duplicate keys, missing fields and nulls, which normal
// struct decoding otherwise silently accepts. The endpoint intentionally has two fields.
func decodeInput(body io.Reader) (service.Input, error) {
	d := json.NewDecoder(body)
	var in service.Input
	tok, err := d.Token()
	if err != nil {
		return in, err
	}
	if tok != json.Delim('{') {
		return in, errors.New("expected object")
	}
	seen := map[string]bool{}
	for d.More() {
		t, err := d.Token()
		if err != nil {
			return in, err
		}
		k, ok := t.(string)
		if !ok || seen[k] {
			return in, errors.New("duplicate key")
		}
		seen[k] = true
		var raw json.RawMessage
		if err := d.Decode(&raw); err != nil {
			return in, err
		}
		if string(raw) == "null" {
			return in, errors.New("null field")
		}
		switch k {
		case "service_id":
			err = json.Unmarshal(raw, &in.ServiceID)
		case "integrations":
			err = json.Unmarshal(raw, &in.Integrations)
		default:
			return in, errors.New("unknown field")
		}
		if err != nil {
			return in, err
		}
	}
	if _, err := d.Token(); err != nil {
		return in, err
	}
	if !seen["service_id"] || !seen["integrations"] {
		return in, errors.New("missing field")
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		if err == nil {
			err = errors.New("trailing JSON")
		}
		return in, err
	}
	return in, nil
}
