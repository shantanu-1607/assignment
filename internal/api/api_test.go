package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shantanu-1607/assignment/internal/service"
)

func request(a http.Handler, method, path, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	return w
}
func newAPI() *API { return New(service.Calculator{}, Options{Rate: 10000}) }
func assertError(t *testing.T, w *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("status=%d want=%d body=%s", w.Code, status, w.Body)
	}
	var e envelope
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if e.Error == nil || e.Error.Code != code || e.RequestID == "" || e.RequestID != w.Header().Get("X-Request-ID") {
		t.Fatalf("inconsistent error: %s", w.Body)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Fatal("error content type is not JSON")
	}
	if w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("error was cacheable")
	}
}
func TestEstimateValidation(t *testing.T) {
	cases := []struct {
		name, body string
		status     int
		code       string
	}{
		{"success", `{"service_id":"api-review","integrations":2}`, 200, ""},
		{"lower", `{"service_id":"api-review","integrations":0}`, 200, ""},
		{"upper", `{"service_id":"api-review","integrations":20}`, 200, ""},
		{"negative", `{"service_id":"api-review","integrations":-1}`, 422, "validation_failed"},
		{"over", `{"service_id":"api-review","integrations":21}`, 422, "validation_failed"},
		{"unknown service", `{"service_id":"secret","integrations":1}`, 422, "validation_failed"},
		{"unknown key", `{"service_id":"api-review","integrations":1,"admin":true}`, 400, "invalid_json"},
		{"duplicate", `{"service_id":"api-review","integrations":1,"integrations":2}`, 400, "invalid_json"},
		{"missing", `{"service_id":"api-review"}`, 400, "invalid_json"},
		{"null", `{"service_id":"api-review","integrations":null}`, 400, "invalid_json"},
		{"null object", `null`, 400, "invalid_json"},
		{"array", `[]`, 400, "invalid_json"},
		{"float", `{"service_id":"api-review","integrations":1.2}`, 400, "invalid_json"},
		{"trailing", `{"service_id":"api-review","integrations":1} {}`, 400, "invalid_json"},
		{"malformed", `{"service_id":`, 400, "invalid_json"},
		{"oversized", `{"service_id":"` + strings.Repeat("x", 4096) + `","integrations":1}`, 413, "body_too_large"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := request(newAPI(), "POST", "/v1/estimates", tc.body)
			if tc.code != "" {
				assertError(t, w, tc.status, tc.code)
				return
			}
			if w.Code != 200 {
				t.Fatalf("%d %s", w.Code, w.Body)
			}
			var result struct {
				Data      service.Estimate `json:"data"`
				RequestID string           `json:"request_id"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			var in service.Input
			_ = json.Unmarshal([]byte(tc.body), &in)
			if result.Data.EffortUnits != 3+2*in.Integrations || result.RequestID == "" {
				t.Fatalf("%s", w.Body)
			}
		})
	}
}
func TestBodyLimitBoundary(t *testing.T) {
	body := `{"service_id":"api-review","integrations":0}`
	for _, size := range []int{4096, 4097} {
		w := request(newAPI(), "POST", "/v1/estimates", body+strings.Repeat(" ", size-len(body)))
		if size == 4096 && w.Code != 200 {
			t.Fatalf("exact limit: %s", w.Body)
		}
		if size == 4097 {
			assertError(t, w, 413, "body_too_large")
		}
	}
}
func TestRoutingAndHeaders(t *testing.T) {
	a := newAPI()
	assertError(t, request(a, "GET", "/missing", ""), 404, "not_found")
	w := request(a, "POST", "/v1/services", "")
	assertError(t, w, 405, "method_not_allowed")
	if w.Header().Get("Allow") != "GET, HEAD" {
		t.Fatal("missing Allow")
	}
	w = request(a, "GET", "/v1/estimates", "")
	assertError(t, w, 405, "method_not_allowed")
	r := httptest.NewRequest("POST", "/v1/estimates", strings.NewReader("{}"))
	r.Header.Set("Content-Type", "text/plain")
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assertError(t, w, 415, "unsupported_media_type")
	w = request(a, "GET", "/healthz", "")
	for _, key := range []string{"Content-Security-Policy", "X-Content-Type-Options", "X-Frame-Options", "Referrer-Policy", "X-Request-ID"} {
		if w.Header().Get(key) == "" {
			t.Fatal(key)
		}
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "" || len(w.Result().Cookies()) > 0 {
		t.Fatal("unexpected browser authority")
	}
	a.SetReady(false)
	assertError(t, request(a, "GET", "/readyz", ""), 503, "not_ready")
	assertError(t, request(a, "GET", "/v1/services", ""), 503, "not_ready")
	if request(a, "GET", "/healthz", "").Code != 200 {
		t.Fatal("liveness failed while draining")
	}
}
func TestCatalogCaching(t *testing.T) {
	a := newAPI()
	w := request(a, "GET", "/v1/services", "")
	var payload struct {
		Data []service.Offering `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil || len(payload.Data) != 2 {
		t.Fatalf("%s", w.Body)
	}
	etag := w.Header().Get("ETag")
	if etag == "" || w.Header().Get("Cache-Control") != "public, max-age=60, must-revalidate" {
		t.Fatal("missing cache policy")
	}
	for _, tag := range []string{etag, "W/" + etag, `"other", ` + etag, "*"} {
		r := httptest.NewRequest("GET", "/v1/services", nil)
		r.Header.Set("If-None-Match", tag)
		out := httptest.NewRecorder()
		a.ServeHTTP(out, r)
		if out.Code != 304 || out.Body.Len() != 0 {
			t.Fatalf("%q: %d", tag, out.Code)
		}
	}
	head := request(a, "HEAD", "/v1/services", "")
	if head.Code != 200 || head.Body.Len() != 0 || head.Header().Get("ETag") != etag {
		t.Fatal("HEAD mismatch")
	}
}
func TestRateLimitIgnoresForwardedHeaders(t *testing.T) {
	a := New(service.Calculator{}, Options{Rate: 1})
	if request(a, "GET", "/v1/services", "").Code != 200 {
		t.Fatal("first request rejected")
	}
	r := httptest.NewRequest("GET", "/v1/services", nil)
	r.Header.Set("X-Forwarded-For", "1.2.3.4")
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assertError(t, w, 429, "rate_limited")
	if w.Header().Get("Retry-After") != "60" {
		t.Fatal("missing retry hint")
	}
	if request(a, "GET", "/readyz", "").Code != 200 {
		t.Fatal("probe subject to quota")
	}
}

type estimatorFunc func(context.Context, service.Input) (service.Estimate, error)

func (f estimatorFunc) Estimate(ctx context.Context, in service.Input) (service.Estimate, error) {
	return f(ctx, in)
}
func TestFailureHygiene(t *testing.T) {
	for _, panicMode := range []bool{false, true} {
		svc := estimatorFunc(func(context.Context, service.Input) (service.Estimate, error) {
			if panicMode {
				panic("secret database password")
			}
			return service.Estimate{}, errors.New("secret database password")
		})
		w := request(New(svc, Options{}), "POST", "/v1/estimates", `{"service_id":"api-review","integrations":1}`)
		assertError(t, w, 500, "internal_error")
		if strings.Contains(w.Body.String(), "password") {
			t.Fatal("leaked internal error")
		}
	}
}
func TestDeadlineCancelsWorker(t *testing.T) {
	canceled := make(chan struct{})
	svc := estimatorFunc(func(ctx context.Context, _ service.Input) (service.Estimate, error) {
		<-ctx.Done()
		close(canceled)
		return service.Estimate{}, ctx.Err()
	})
	a := New(svc, Options{Timeout: 10 * time.Millisecond})
	w := request(a, "POST", "/v1/estimates", `{"service_id":"api-review","integrations":1}`)
	assertError(t, w, 503, "timeout")
	select {
	case <-canceled:
	case <-time.After(time.Second):
		t.Fatal("worker not canceled")
	}
}
func TestClientCancellation(t *testing.T) {
	canceled := make(chan struct{})
	svc := estimatorFunc(func(ctx context.Context, _ service.Input) (service.Estimate, error) {
		<-ctx.Done()
		close(canceled)
		return service.Estimate{}, ctx.Err()
	})
	a := New(svc, Options{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := httptest.NewRequest("POST", "/v1/estimates", strings.NewReader(`{"service_id":"api-review","integrations":1}`)).WithContext(ctx)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	if w.Code != 503 {
		t.Fatalf("status=%d", w.Code)
	}
	select {
	case <-canceled:
	case <-time.After(time.Second):
		t.Fatal("client cancellation lost")
	}
}
func TestAdmissionBound(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	svc := estimatorFunc(func(context.Context, service.Input) (service.Estimate, error) {
		close(entered)
		<-release
		return service.Estimate{}, nil
	})
	a := New(svc, Options{MaxConcurrent: 1})
	done := make(chan struct{})
	go func() {
		defer close(done)
		request(a, "POST", "/v1/estimates", `{"service_id":"api-review","integrations":1}`)
	}()
	<-entered
	w := request(a, "GET", "/v1/services", "")
	assertError(t, w, 503, "busy")
	close(release)
	<-done
	if request(a, "GET", "/v1/services", "").Code != 200 {
		t.Fatal("slot not released")
	}
}
func TestTimeoutRetainsSlotUntilWorkerStops(t *testing.T) {
	release, finished := make(chan struct{}), make(chan struct{})
	svc := estimatorFunc(func(context.Context, service.Input) (service.Estimate, error) {
		<-release
		defer close(finished)
		return service.Estimate{}, nil
	})
	a := New(svc, Options{Timeout: 5 * time.Millisecond, MaxConcurrent: 1})
	assertError(t, request(a, "POST", "/v1/estimates", `{"service_id":"api-review","integrations":1}`), 503, "timeout")
	assertError(t, request(a, "GET", "/v1/services", ""), 503, "busy")
	close(release)
	<-finished
}
func TestLimiterBoundAndExpiry(t *testing.T) {
	now := time.Unix(0, 0)
	l := &limiter{entries: map[string]bucket{}, limit: 2, maxKeys: 1, now: func() time.Time { return now }}
	if !l.allow("a") || !l.allow("a") || l.allow("a") || l.allow("b") {
		t.Fatal("limit or capacity incorrect")
	}
	now = now.Add(time.Minute)
	if !l.allow("b") || len(l.entries) != 1 {
		t.Fatal("expired entries not reclaimed")
	}
}
func TestLimiterConcurrent(t *testing.T) {
	l := &limiter{entries: map[string]bucket{}, limit: 17, maxKeys: 10, now: time.Now}
	var mu sync.Mutex
	allowed := 0
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.allow("same") {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if allowed != 17 {
		t.Fatalf("allowed=%d", allowed)
	}
}
func TestLiveHTTP(t *testing.T) {
	s := httptest.NewServer(newAPI())
	defer s.Close()
	resp, err := s.Client().Get(s.URL + "/v1/services")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 || !strings.Contains(string(body), "api-review") {
		t.Fatalf("%s %v", body, err)
	}
}
