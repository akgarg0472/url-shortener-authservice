package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"

	DB "github.com/akgarg0472/urlshortener-auth-service/database"
	DiscoveryClient "github.com/akgarg0472/urlshortener-auth-service/discovery-client"
	Routers "github.com/akgarg0472/urlshortener-auth-service/internal/router"
	OAuthService "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth/oauth"
	KafkaService "github.com/akgarg0472/urlshortener-auth-service/internal/service/kafka"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	Utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

func init() {
	loadDotEnv()
	DB.InitDB()
	OAuthService.InitOAuthProviders()
	KafkaService.InitKafka()
}

var (
	logger    = Logger.GetLogger("main.go")
	BuildTime string
	GitCommit string
	BuildHost string
	BuildEnv  string
)

func main() {
	// Set up a context to manage the server's shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a channel to capture termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	portEnv := Utils.GetEnvVariable("SERVER_PORT", "8081")
	port, err := strconv.Atoi(portEnv)

	if err != nil {
		panic("Invalid port value defined in environment: " + portEnv)
	}

	DiscoveryClient.InitDiscoveryClient(port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: loadRoutersV1(),
	}

	go func() {
		logger.Info("Starting server on port: {}", port)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Error starting server: {}", err)
		} else {
			logger.Info("Server started on port: {}", port)
		}
	}()

	defer cleanupResources(server)

	// Wait for a termination signal
	<-sigCh

	logger.Info("Shutting down server gracefully...")

	shutdownError := server.Shutdown(ctx)

	if shutdownError != nil {
		logger.Error("Error during server shutdown: {}", err)
	}
}

func loadDotEnv() {
	err := godotenv.Load()

	if err != nil {
		logger.Info("Error loading .env file: %v", err)
	}
}

func loadRoutersV1() *chi.Mux {
	router := chi.NewRouter()

	router.Mount("/api/v1/auth", Routers.AuthRouterV1())
	router.Mount("/api/v1/auth/oauth", Routers.OAuthRouterV1())
	router.Mount("/", Routers.PingRouterV1())
	router.Mount("/admin", Routers.DiscoveryRouterV1())

	return router
}

func cleanupResources(server *http.Server) {
	logger.Info("Cleaning up before exiting...")

	dbCloseError := DB.CloseDB()

	if dbCloseError != nil {
		logger.Error("Error closing DB connection: {}", dbCloseError.Error())
	}

	discoveryClientCloseError := DiscoveryClient.UnregisterInstance()

	if discoveryClientCloseError != nil {
		logger.Error("Error unregistering discovery Client: {}", discoveryClientCloseError.Error())
	}

	kafkaCloseError := KafkaService.CloseKafka()

	if kafkaCloseError != nil {
		logger.Error("Error closing kafka connection: {}", kafkaCloseError.Error())
	}

	if server != nil {
		serverCloseError := server.Shutdown(context.Background())
		if serverCloseError != nil {
			logger.Error("Error shutting down server: {}", serverCloseError.Error())
		}
	}
}
