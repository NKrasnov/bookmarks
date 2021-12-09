// web implements simple api framework
package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler func(context.Context, http.ResponseWriter, *http.Request) error

type Middleware func(Handler) Handler

// App
type App struct {
	mux.Router
	shutdown <-chan os.Signal
	mw       []Middleware
	log      *log.Logger
}

func NewApp(s <-chan os.Signal, log *log.Logger, m ...Middleware) *App {
	return &App{
		shutdown: s,
		mw:       m,
		log:      log,
	}
}

type ctxKey int

const ContextKeyValue ctxKey = 1

type ContextValue struct {
	TraceID string
	Now     time.Time
	Status  int
}

func (a *App) Handle(path string, handler Handler, mw ...Middleware) {
	// wrap handler into specific handler middlewares
	handler = wrapMiddleware(mw, handler)
	// then wrap handler into general App middlewares
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		a.log.Println("App handle called. cerated context and passed down the request pipeline")

		v := ContextValue{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), ContextKeyValue, &v)
		if err := handler(ctx, w, r); err != nil {
			a.log.Println(err)
		}
	}
	a.HandleFunc(path, h)
}

func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}
	return handler
}

func GetValueFromContext(ctx context.Context) (*ContextValue, error) {
	v, ok := ctx.Value(ContextKeyValue).(*ContextValue)
	if !ok {
		return nil, fmt.Errorf("cannot retrieve value from context")
	}
	return v, nil
}
