package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/akgarg0472/urlshortener-auth-service/database"
	"github.com/akgarg0472/urlshortener-auth-service/discovery"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	"github.com/akgarg0472/urlshortener-auth-service/internal/metrics"
	"github.com/akgarg0472/urlshortener-auth-service/internal/router"
	oauth_service "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth/oauth"
	kafka_service "github.com/akgarg0472/urlshortener-auth-service/internal/service/kafka"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	BuildTime string
	GitCommit string
	BuildHost string
	BuildEnv  string
)

func init() {
	database.InitDB()
	oauth_service.InitOAuthProviders()
	kafka_service.InitKafka()
}

func main() {
	// Set up a context to manage the server's shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a channel to capture termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	portEnv := utils.GetEnvVariable("SERVER_PORT", "8081")
	port, err := strconv.Atoi(portEnv)

	if err != nil {
		panic(fmt.Sprintf("Invalid port value defined in environment: %s", portEnv))
	}

	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		logger.Error("Error creating TCP listener", zap.Error(err))
		return
	}

	actualPort, err := extractActualServerPort(listener.Addr().String())

	if err != nil {
		logger.Error("Error resolving host and port", zap.Error(err))
		return
	}

	logger.AddPortToLogger(actualPort)

	server := &http.Server{
		Handler: loadRoutersV1(),
	}

	discovery.InitDiscoveryClient(actualPort)

	go func() {
		err := server.Serve(listener)

		if err != nil {
			logger.Error("Error starting server", zap.Error(err))
		}
	}()

	defer cleanupResources(server)

	// Wait for a termination signal
	<-sigCh

	if logger.IsInfoEnabled() {
		logger.Info("Shutting down server gracefully...")
	}

	shutdownError := server.Shutdown(ctx)

	if shutdownError != nil {
		logger.Error("Error during server shutdown", zap.Error(err))
	}
}

func loadRoutersV1() *chi.Mux {
	r := chi.NewRouter()

	r.Use(metrics.PrometheusMiddleware)

	r.Mount("/api/v1/auth", router.AuthRouterV1())
	r.Mount("/api/v1/auth/oauth", router.OAuthRouterV1())
	r.Mount("/", router.PingRouterV1())
	r.Mount("/admin", router.DiscoveryRouterV1())
	r.Handle("/prometheus/metrics", metrics.MetricsHandler())

	return r
}

func cleanupResources(server *http.Server) {
	if logger.IsInfoEnabled() {
		logger.Info("Cleaning up before exiting...")
	}

	if err := database.CloseDB(); err != nil && logger.IsErrorEnabled() {
		if logger.IsErrorEnabled() {
			logger.Error("Error closing DB connection", zap.Error(err))
		}
	}

	if err := discovery.UnregisterInstance(); err != nil && logger.IsErrorEnabled() {
		if logger.IsErrorEnabled() {
			logger.Error("Error unregistering discovery client", zap.Error(err))
		}
	}

	if err := kafka_service.CloseKafka(); err != nil && logger.IsErrorEnabled() {
		if logger.IsErrorEnabled() {
			logger.Error("Error closing Kafka connection", zap.Error(err))
		}
	}

	if server != nil {
		if err := server.Shutdown(context.Background()); err != nil && logger.IsErrorEnabled() {
			if logger.IsErrorEnabled() {
				logger.Error("Error shutting down server", zap.Error(err))
			}
		}
	}
}

func extractActualServerPort(addr string) (int, error) {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(portStr)
}
