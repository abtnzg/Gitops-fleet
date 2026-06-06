// service-a: minimal HTTP service used to validate the platform end-to-end.
// Exposes /, /healthz, /version, /echo. Logs to stdout in JSON. No deps beyond
// the standard library so the build is hermetic and the binary is small.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	// Set at build time via -ldflags "-X main.version=...".
	version = "dev"
	// Set at build time via -ldflags "-X main.commit=...".
	commit = "unknown"
	// Set at build time via -ldflags "-X main.buildTime=...".
	buildTime = "unknown"
)

type response struct {
	Service   string `json:"service"`
	Message   string `json:"message,omitempty"`
	Echo      string `json:"echo,omitempty"`
	Method    string `json:"method,omitempty"`
	Path      string `json:"path,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	Timestamp string `json:"timestamp"`
	Hostname  string `json:"hostname"`
}

func main() {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/healthz", handleHealthz)
	mux.HandleFunc("/version", handleVersion)
	mux.HandleFunc("/echo", handleEcho)

	// Log startup so it's visible in `kubectl logs` immediately.
	log.Printf("service-a starting version=%s commit=%s built=%s go=%s addr=%s",
		version, commit, buildTime, runtime.Version(), addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      logRequests(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, response{
		Service:   "service-a",
		Message:   "hello from the platform-fleet sandbox cluster",
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Hostname:  hostname(),
	})
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, response{
		Service:   "service-a",
		Message:   "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Hostname:  hostname(),
	})
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, struct {
		Service   string `json:"service"`
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		BuildTime string `json:"build_time"`
		GoVersion string `json:"go_version"`
		Hostname  string `json:"hostname"`
	}{
		Service:   "service-a",
		Version:   version,
		Commit:    commit,
		BuildTime: buildTime,
		GoVersion: runtime.Version(),
		Hostname:  hostname(),
	})
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	// Limit echoed body to 1 KB to keep responses small for the test.
	const maxEcho = 1024
	if len(body) > maxEcho {
		body = body[:maxEcho]
	}
	headers := map[string]string{}
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	writeJSON(w, http.StatusOK, response{
		Service:   "service-a",
		Echo:      string(body),
		Method:    r.Method,
		Path:      r.URL.Path,
		Headers:   headers,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Hostname:  hostname(),
	})
}

func logRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		h.ServeHTTP(rw, r)
		log.Printf("method=%s path=%s status=%d duration_ms=%d bytes=%d",
			r.Method, r.URL.Path, rw.status, time.Since(start).Milliseconds(), rw.bytes)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(s int) {
	r.status = s
	r.ResponseWriter.WriteHeader(s)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "encode response: %v\n", err)
	}
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}
