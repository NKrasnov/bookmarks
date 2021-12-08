package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nkrasnov/bookmarks/server/foundation/config"
)

func main() {
	log := log.New(os.Stdout, "BM_API: ", log.Lshortfile)
	err := run(os.Args, log)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
}

func run(args []string, log *log.Logger) error {

	// App configuration struct
	cfg := struct {
		APIHost               string        `param:"cmd=host,env=BM_API_HOST,default=127.0.0.1,usage=IP or DNS Name"`
		APIPort               int           `param:"cmd=port,env=BM_API_PORT,default=8081,usage=API server port"`
		DBHost                string        `param:"cmd=dbhost,env=BM_API_DBHOST,default=127.0.0.1,usage=IP or DNS Name"`
		DBPort                int           `param:"cmd=dbport,env=BM_API_DBPORT,default=5432,usage=Port number a database server is listens to"`
		DBUser                string        `param:"cmd=dbuser,env=BM_API_DBUSER,default=postgres,usage=database user name"`
		DBPwd                 string        `param:"cmd=dbpwd,env=BM_API_DBPWD,usage=database user password"`
		ReadTimeout           time.Duration `param:"cmd=srto,env=BM_API_READ_TIMEOUT,default=10,usage=API server read timeout"`
		WriteTimeout          time.Duration `param:"cmd=swto,env=BM_API_WRITE_TIMEOUT,default=10,usage=API server write timeout"`
		RequestTimeout        time.Duration `param:"cmd=rto,env=BM_API_REQUEST_TIMEOUT,default=5,usage=API request timeout"`
		ServerShutdownTimeout time.Duration `param:"cmd=ssto,env=BM_API_SHUTDOWN_TIMEOUT,default=30,usage=Time during which server allowed to finish all it's work after shutdown has been initiated"`
	}{}
	err := config.Parse(&cfg, args)
	if errors.Is(err, config.ErrHelpNeeded) {
		config.PrintUsage()
		return nil
	}

	// API Server setup
	shutdown := make(chan os.Signal, 1)
	apiServerError := make(chan error, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.APIHost, cfg.APIPort),
		ReadTimeout:  cfg.ReadTimeout * time.Second,
		WriteTimeout: cfg.WriteTimeout * time.Second,
		Handler:      http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Hello there") })),
	}
	go func() {
		log.Printf("server is listening on %s:%d address\n", cfg.APIHost, cfg.APIPort)
		apiServerError <- api.ListenAndServe()
	}()

	select {
	case sig := <-shutdown:
		log.Printf("main: Shutdown signal received: %v", sig)
		log.Println("main: Waiting for requests to complete...")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ServerShutdownTimeout*time.Second)
		defer cancel()
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could'n shutdown API server %w", err)
		}
	case err := <-apiServerError:
		return fmt.Errorf("API server error: %v", err)
	}
	return nil
}
