package main

import (
	"context"
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
	OAuthService.InitOAuthClients()
	KafkaService.InitKafka()
}

var (
	logger = Logger.GetLogger("main.go")
)

func main() {
	// Set up a context to manage the server's shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a channel to capture termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	port, _ := strconv.Atoi(Utils.GetEnvVariable("SERVER_PORT", "8081"))

	DiscoveryClient.InitDiscoveryClient(port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: loadRouters(),
	}

	go func() {
		logger.Info("Starting server on port: {}", port)

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("Error starting server: {}", err)
		} else {
			logger.Info("Server started on port: {}", port)
		}
	}()

	defer cleanupResources(server)

	// Wait for a termination signal
	<-sigCh

	// Start the graceful server shutdown
	logger.Info("Shutting down server gracefully...")
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Error("Error during server shutdown: {}", err)
	}
}

func loadDotEnv() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}
}

func loadRouters() *chi.Mux {
	router := chi.NewRouter()

	router.Mount("/auth/v1", Routers.AuthRouterV1())
	router.Mount("/auth/v1/oauth", Routers.OAuthRouterV1())

	return router
}

func cleanupResources(server *http.Server) {
	logger.Info("Cleaning up before exiting...")

	dbCloseError := DB.CloseDB()

	if dbCloseError != nil {
		logger.Error("Error closing DB connection: {}", dbCloseError.Error())
	}

	discovertClientCloseError := DiscoveryClient.UnregisterInstance()

	if discovertClientCloseError != nil {
		logger.Error("Error unregistering discovery Client: {}", discovertClientCloseError.Error())
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
