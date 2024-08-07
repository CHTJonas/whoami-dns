package whoami

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

type Server struct {
	bin *cache.Cache
	web *http.Server
	in  *dnstap.FrameStreamSockInput
	out *dnstap.TextOutput
	mwf []mux.MiddlewareFunc
}

func NewServer() *Server {
	return &Server{
		bin: cache.New(30*time.Second, 2*time.Minute),
		mwf: []mux.MiddlewareFunc{headerMiddleware},
	}
}

func (s *Server) Write(p []byte) (n int, err error) {
	str := string(p)
	parts := strings.Split(str, " ")
	ip := parts[2]
	hostname := parts[5]
	hostname = strings.Trim(hostname, `"`)
	hostname = strings.TrimSuffix(hostname, ".")
	hostname = strings.ToLower(hostname)
	fmt.Println("DNS query for", hostname, "from", ip)
	s.bin.SetDefault(hostname, ip)
	return len(p), nil
}

func (s *Server) whoamiEndpoint(w http.ResponseWriter, r *http.Request) {
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
	h = strings.ToLower(h)
	fmt.Println("HTTP request for", h, "from", ip)

	for i := 0; i < 20; i++ {
		body, found := s.bin.Get(h)
		if found {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body.(string) + "\n"))
			return
		}
		time.Sleep(500 * time.Millisecond)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Error: no query received\n"))
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, max-age=0")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) SetHeader(name, value string) {
	s.mwf = append(s.mwf, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(name, value)
			next.ServeHTTP(w, r)
		})
	})
}

func (s *Server) OpenSocket(path string) {
	input, err := dnstap.NewFrameStreamSockInputFromPath(path)
	if err != nil {
		panic(err)
	}
	output := dnstap.NewTextOutput(s, dnstap.TextFormat)
	s.in = input
	s.out = output
	go output.RunOutputLoop()
	go input.ReadInto(output.GetOutputChannel())
	fmt.Println("dnstap socket opened at", path)
}

func (s *Server) CloseSocket() {
	s.out.Close()
}

func (s *Server) Start(port string) {
	router := mux.NewRouter()
	router.HandleFunc("/", s.whoamiEndpoint).Methods("GET", "HEAD")
	router.Use(s.mwf...)

	s.web = &http.Server{
		Addr:        ":" + port,
		Handler:     router,
		ReadTimeout: 30 * time.Second,
	}
	s.web.SetKeepAlivesEnabled(false)
	go func() {
		if err := s.web.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP server error:", err)
			os.Exit(99)
		}
	}()
	fmt.Println("HTTP server listening on", s.web.Addr)
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fmt.Println("Waiting up to 30 seconds for HTTP server to shutdown")
	if err := s.web.Shutdown(ctx); err != nil {
		fmt.Println("HTTP server shutdown error:", err)
	}
	fmt.Println("HTTP server gone. Ta-ta!")
}
