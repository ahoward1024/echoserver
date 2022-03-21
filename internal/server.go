package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type httpServer interface {
	Shutdown(ctx context.Context) error
	ListenAndServe() error
}

type Server struct {
	http.Server
	Name string
}

type Options struct {
	LivenessFilePath string
	Host             string
	Port             int
	MetricsPort      int
	Wait             int
	WriteTimeout     int
	ReadTimeout      int
	IdleTimeout      int
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg(fmt.Sprintf("Started request to %s", r.RequestURI))
		next.ServeHTTP(w, r)
		log.Debug().Msg(fmt.Sprintf("Finished request to %s", r.RequestURI))
	})
}

func serverShutdown(server *Server, ctx context.Context) {
	log.Info().Msg(fmt.Sprintf("%s server shutting down", server.Name))
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Msg(fmt.Sprintf("%s server shutdown unsuccessful: %v", server.Name, err))
	}
	log.Info().Msg(fmt.Sprintf("%s server shut down successfully", server.Name))
}

func serverStartup(server *Server) {
	log.Info().Msg(fmt.Sprintf("%s server startup", server.Name))
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Msg(fmt.Sprintf("%s server startup unsuccessful: %v", server.Name, err))
	}
}

func RunServer(opts *Options) error {
	// Setup the liveness file
	path := opts.LivenessFilePath
	livenessFile, err := os.Create(path)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Failed to create liveness file: %s", err))
		return fmt.Errorf("Failed to create liveness file: %s", err)
	}

	// Setup the channel for blocking
	defer os.Remove(livenessFile.Name())
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(opts.Wait))
	defer cancel()

	// Setup the main server
	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	AddRoutes(router)
	server := &Server{
		http.Server{
			Addr:         fmt.Sprintf("%s:%d", opts.Host, opts.Port),
			WriteTimeout: time.Second * time.Duration(opts.WriteTimeout),
			ReadTimeout:  time.Second * time.Duration(opts.ReadTimeout),
			IdleTimeout:  time.Second * time.Duration(opts.IdleTimeout),
			Handler:      router,
		},
		"Main",
	}
	defer serverShutdown(server, ctx)

	// If the metrics port is not the same as the server port, run the metrics as a separate server
	if opts.Port != opts.MetricsPort {
		router := mux.NewRouter()
		router.Use(loggingMiddleware)
		AddMetricsRoute(router)
		server := &Server{
			http.Server{
				Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.MetricsPort),
				Handler: router,
			},
			"Metrics",
		}

		// Run the metrics server
		go func() {
			serverStartup(server)
		}()
		defer serverShutdown(server, ctx)
		loggingMiddleware(router)
	} else {
		// Otherwise we will add the /metrics route to the main server
		AddMetricsRoute(router)
	}

	// Run the main server
	go func() {
		serverStartup(server)
	}()

	<-shutdown // Block until shutdown

	return err
}
