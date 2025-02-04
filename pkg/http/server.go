package http

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/renatospaka/customer-processor-test/pkg/http/ip"
	middlewares "github.com/renatospaka/customer-processor-test/pkg/http/middleware"
	"github.com/rs/zerolog/log"
)

var (
	ErrHTTPNilHandler = fmt.Errorf("http: nil handler, could not add route to HTTP server")
	ErrHTTPNilPrefix  = fmt.Errorf("nil prefix, could not create/find a group with nil prefix")
	ErrHTTPNilMethod  = fmt.Errorf("nil method, could not create/find a route with nil method")
	ErrHTTPNilGroup   = fmt.Errorf("there is no routing group defined")
	ErrHTTPNilRoute   = fmt.Errorf("there is no route defined on this server")
	ErrHTTPNilRouter  = fmt.Errorf("there is no router defined on this server")

	ErrHTTPServerAddressNotProvided = fmt.Errorf("server address not provided ")
	ErrHTTPServerNotInitialized     = fmt.Errorf("HTTP server not initialized")
)

type HTTPServer struct {
	mu     sync.RWMutex
	server *http.Server
	mux    *chi.Mux

	addr    string
	port    int
	version string
}

// New returns a fresh instance of HTTPServer from basic attributes set
func New(opts ...Options) *HTTPServer {
	// Initialize build options
	httpServer := &HTTPServer{
		mu:  sync.RWMutex{},
		mux: &chi.Mux{},
	}
	for _, opt := range opts {
		opt(httpServer)
	}

	// Verify if address is a valid IP
	// Return nil if doesn't
	if err := ip.IsValidIP(httpServer.addr); err != nil {
		log.Error().Msgf("error using address (%s): %s",
			httpServer.addr,
			err.Error())
		return nil
	}

	// Initialize HTTP Server & MUX
	address := ip.FormatAddress(httpServer.addr, httpServer.port)
	httpServer.addr = address
	httpServer.mux = newRouter().mux

	// Configure HTTP
	server := &http.Server{
		Addr:    address,
		Handler: httpServer.mux,
	}
	httpServer.server = server
	return httpServer
}

// Addr is the address where server is listening
func (h *HTTPServer) Addr() string {
	return h.addr
}

// Version of the application running http
func (h *HTTPServer) Version() string {
	return h.version
}

// Server is the server, listening or not
func (h *HTTPServer) Server() *http.Server {
	return h.server
}

// Mux is the multiplexer
func (h *HTTPServer) Mux() *chi.Mux {
	return h.mux
}

type Options func(*HTTPServer)

// Informs the address to listen to
func WithAddress(addr string) Options {
	return func(h *HTTPServer) {
		h.addr = addr
	}
}

// Informs the port of the app, if any
func WithPort(port int) Options {
	return func(h *HTTPServer) {
		h.port = port
	}
}

// Informs the version of the app, if any
func WithVersion(version string) Options {
	return func(h *HTTPServer) {
		h.version = version
	}
}

// HTTPServer receives address and an instance of *http.Server
// and initializes the HTTP server on the port read by the service setup (*Setup)
func (h *HTTPServer) Serve(ready chan<- bool) (*http.Server, error) {
	if h == nil {
		return nil, ErrHTTPServerNotInitialized
	}

	address := h.addr
	log.Info().Msg("starting HTTP server")

	var err error
	go func() {
		// sign that server is ready to server
		ready <- true

		if err = h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("error when server started listening: %v", err)
			h.server = nil
		}
	}()

	if h.server == nil {
		return nil, ErrHTTPServerNotInitialized
	}

	log.Info().Msgf("HTTP server listening on %s", address)
	return h.server, err
}

// ServeWithTLS receives address and an instance of *http.Server with TLS features
// and initializes the HTTP server on the port read by the service setup (*Setup)
func (h *HTTPServer) ServeWithTLS(ready chan<- bool) (*http.Server, error) {
	if h == nil {
		return nil, ErrHTTPServerNotInitialized
	}

	log.Warn().Msg("TLS not yet implemented")
	return nil, ErrHTTPServerNotInitialized
}

// Closes the HTTP server
func (h *HTTPServer) Close() error {
	if h == nil {
		return nil
	}

	var erros error
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.server != nil {
		erros = h.server.Close()
		if erros != nil {
			log.Error().Msgf("disconnecting HTTP server with error: %v", erros)
		}
		h.server = nil
	}
	if erros == nil {
		log.Warn().Msgf("successfully disconnected from HTTP server")
	}
	return erros
}

type router struct {
	mux *chi.Mux
}

// Creates a new instance of the router (Chi in this case)
func newRouter() *router {
	m := chi.NewRouter()
	r := &router{
		mux: m,
	}
	r.initLogger()
	r.basicSetup()
	return r
}

// Prepares the http middleware to use Logger (from app) on it
func (r *router) initLogger() {
	r.mux.Use(middlewares.Logger)
}

// Executes basic setup for the router
func (r *router) basicSetup() {
	r.mux.Use(middlewares.Cors)
	r.mux.Use(middleware.Recoverer)
	r.mux.Use(middleware.Heartbeat("/health"))
}

// Returns the multiplexer of the router
func (r *router) Mux() *chi.Mux {
	return r.mux
}
