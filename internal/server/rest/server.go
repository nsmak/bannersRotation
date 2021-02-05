package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/nsmak/bannersRotation/internal/app"
)

type serverError struct {
	app.BaseError
}

func newServerError(msg string, err error) *serverError {
	return &serverError{BaseError: app.BaseError{Message: msg, Err: err}}
}

type Server struct {
	Address string
	public  API
	server  *http.Server
	log     app.Logger
}

func NewServer(public API, address string, logger app.Logger) *Server {
	return &Server{
		Address: address,
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
		return newServerError("start server error", err)
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return newServerError("server is nil", nil)
	}

	if err := s.server.Shutdown(ctx); err != nil {
		return newServerError("stop server error", err)
	}
	return nil
}

func (s *Server) router() *mux.Router {
	router := mux.NewRouter()
	for _, r := range s.public.Routes() {
		handler := alice.New(s.panicMiddleware, s.loggingMiddleware).ThenFunc(r.Func)
		router.
			Methods(r.Method).
			Path(r.Path).
			Name(r.Name).
			Handler(handler)
	}
	return router
}
