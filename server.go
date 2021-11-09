package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/gorilla/mux"
)

type server struct {
	store sync.Map
	web   *http.Server
}

func newServer() *server {
	s := &server{}
	return s
}

func (s *server) Write(p []byte) (n int, err error) {
	str := string(p)
	parts := strings.Split(str, " ")
	ip := parts[2]
	hostname := parts[5]
	hostname = strings.Trim(hostname, `"`)
	hostname = strings.TrimSuffix(hostname, ".")
	fmt.Println("DNS query for", hostname, "from", ip)
	s.store.Store(hostname, ip)
	return len(p), nil
}

func (s *server) whoamiEndpoint(w http.ResponseWriter, r *http.Request) {
	ip := r.Header["X-Forwarded-For"][0]
	if ip == "" {
		ip = r.RemoteAddr
	}
	host := r.Header["X-Forwarded-Host"][0]
	if host == "" {
		host = r.Host
	}
	u, _ := url.Parse("http://" + host)
	h := u.Hostname()
	h = strings.TrimSpace(h)
	h = strings.TrimSuffix(h, ".")
	fmt.Println("HTTP request for", h, "from", ip)

	for i := 0; i < 5; i++ {
		body, found := s.store.Load(h)
		if found {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body.(string) + "\n"))
			return
		}
		time.Sleep(2 * time.Second)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Error: no query received\n"))
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "https://github.com/CHTJonas/whoami-dns")
		w.Header().Set("Cache-Control", "no-store, max-age=0")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func (s *server) openSocket() {
	input, err := dnstap.NewFrameStreamSockInputFromPath("/var/lib/knot/dnstap.sock")
	if err != nil {
		panic(err)
	}
	output := dnstap.NewTextOutput(s, dnstap.TextFormat)
	go output.RunOutputLoop()
	input.ReadInto(output.GetOutputChannel())
}

func (s *server) start() {
	router := mux.NewRouter()
	router.HandleFunc("/", s.whoamiEndpoint).Methods("GET")
	router.Use(headerMiddleware)

	s.web = &http.Server{Addr: ":6780", Handler: router}
	go func() {
		if err := s.web.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP server error:", err)
			os.Exit(99)
		}
	}()
	fmt.Println("HTTP server listening on", s.web.Addr)
}

func (s *server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fmt.Println("Waiting up to 30 seconds for HTTP server to shutdown")
	if err := s.web.Shutdown(ctx); err != nil {
		fmt.Println("HTTP server shutdown error:", err)
	}
	fmt.Println("HTTP server gone. Ta-ta!")
}
