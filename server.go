package gateway

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/renevo/gateway/logging"
	"github.com/renevo/gateway/static"
)

// Server represents a gateway server instance
type Server struct {
	inner *http.Server
	mux   *http.ServeMux // TODO: a better mux
	site  *static.Site
}

// New creates a new server instance
func New(options ...Option) *Server {
	server := &Server{
		mux:  http.NewServeMux(),
		site: static.New("./public/www"),
	}

	server.inner = &http.Server{
		Handler:           server,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       time.Second,
		WriteTimeout:      time.Second,
	}

	for _, opt := range options {
		opt(server)
	}

	server.mux.Handle("/", server.site)

	return server
}

// Listen will create a new listener and serve requests on it
func (s *Server) Listen(addr *url.URL) error {
	network := addr.Scheme
	if network == "" {
		network = "tcp"
	}

	address := addr.Hostname()
	port := addr.Port()
	if port == "" {
		port = "80"
	}

	ln, err := net.Listen(network, address+":"+port)
	if err != nil {
		return err
	}

	logging.Infof("Serving HTTP requests on %s", ln.Addr())
	return s.inner.Serve(tcpKeepAliveListener{ln.(*net.TCPListener), s.inner.IdleTimeout})
}

// ServeHTTP is the core HTTP handler for the gateway
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	stats := &responseWriterStats{inner: w}
	s.mux.ServeHTTP(stats, r)
	logging.Infof("HTTP %s %s %q %q %s %d %d", r.Method, r.RemoteAddr, r.RequestURI, r.UserAgent(), time.Since(start), stats.code, stats.size)
}

// Shutdown will gracefully shutdown the server, finishing any finalized requests
func (s *Server) Shutdown(ctx context.Context) error {
	return s.inner.Shutdown(ctx)
}

type responseWriterStats struct {
	inner http.ResponseWriter
	code  int
	size  int
}

func (r *responseWriterStats) Header() http.Header {
	return r.inner.Header()
}

func (r *responseWriterStats) Write(v []byte) (int, error) {
	size, err := r.inner.Write(v)
	r.size += size
	return size, err
}

func (r *responseWriterStats) WriteHeader(code int) {
	r.code = code
	r.inner.WriteHeader(code)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
	idleTimeout time.Duration
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(ln.idleTimeout)
	return tc, nil
}
