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
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	DB "github.com/akgarg0472/urlshortener-auth-service/database"
	DiscoveryClient "github.com/akgarg0472/urlshortener-auth-service/discovery-client"
	AuthRouter "github.com/akgarg0472/urlshortener-auth-service/internal/router"
	Utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

func init() {
	loadDotEnv()
	DB.InitDB()
}

func main() {
	// Set up a context to manage the server's shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a channel to capture termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	port := Utils.GetEnvVariable("SERVER_PORT", "8081")
	_port, _ := strconv.Atoi(port)

	DiscoveryClient.InitDiscoveryClient(_port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", _port),
		Handler: loadRouters(),
	}

	go func() {
		fmt.Printf("Starting server on port: %d\n", _port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	defer cleanupResources(server)

	// Wait for a termination signal
	<-sigCh

	// Start the graceful server shutdown
	fmt.Println("Shutting down server gracefully...")
	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
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

	// router.Use(corsHandler())

	router.Mount("/auth/v1", AuthRouter.AuthRouterV1())

	return router
}

func corsHandler() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}

func cleanupResources(server *http.Server) {
	fmt.Println("Cleaning up before exiting...")

	DB.CloseDB()
	DiscoveryClient.UnregisterInstance()

	if server != nil {
		server.Shutdown(context.Background())
	}
}
