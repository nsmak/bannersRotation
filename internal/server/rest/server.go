package rest

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/nsmak/bannersRotation/internal/app"
)

type ServerError struct { // TODO: - переделать!
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

func (e *ServerError) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}
func (e *ServerError) Unwrap() error {
	return e.Err
}

type Server struct {
	Address string
	public  API
	server  *http.Server
	log     app.Logger
}

func NewServer(public API, host, port string, logger app.Logger) *Server {
	return &Server{
		Address: net.JoinHostPort(host, port),
		public:  public,
		log:     logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{
		Addr:         s.Address,
		Handler:      s.router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return &ServerError{Message: "start server error", Err: err}
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return &ServerError{Message: "server is nil"}
	}

	if err := s.server.Shutdown(ctx); err != nil {
		return &ServerError{Message: "stop server error", Err: err}
	}
	return nil
}

func (s *Server) router() *mux.Router {
	router := mux.NewRouter()
	for _, r := range s.public.Routes() {
		handler := alice.New(s.loggingMiddleware).ThenFunc(r.Func)
		router.
			Methods(r.Method).
			Path(r.Path).
			Name(r.Name).
			Handler(handler)
	}
	return router
}
