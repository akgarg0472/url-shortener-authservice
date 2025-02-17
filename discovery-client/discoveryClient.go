package discoveryclient

import (
	"fmt"
	"strconv"
	"time"

	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"github.com/hashicorp/consul/api"
)

var (
	logger       = Logger.GetLogger("discoveryClient.go")
	consulClient *api.Client
	serviceID    string
	serviceName  = "urlshortener-auth-service"
)

func InitDiscoveryClient(port int) {
	isDiscoveryClientEnabled, err := strconv.ParseBool(utils.GetEnvVariable("ENABLE_DISCOVERY_CLIENT", "false"))
	if err != nil || !isDiscoveryClientEnabled {
		logger.Info("Discovery client is disabled in configuration")
		return
	}

	consulAddress := utils.GetEnvVariable("CONSUL_SERVER_IP", "http://127.0.0.1:8500")

	config := api.DefaultConfig()
	config.Address = consulAddress
	consulClient, err = api.NewClient(config)
	if err != nil {
		logger.Fatal("Failed to create Consul client: %v", err)
	}

	host := utils.GetHostIP()
	serviceID = fmt.Sprintf("%s-%d", serviceName, port)

	registerService(port, host)
	initHeartbeat(port, host)
}

func registerService(port int, host string) {
	logger.Info("Registering service with Consul: {}:{}", serviceName, serviceID)

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Port:    port,
		Address: host,
	}

	retryDelay := utils.GetEnvDurationSeconds("REGISTER_RETRY_DELAY_SECONDS", 5*time.Second)
	maxRetryDuration := utils.GetEnvDurationSeconds("REGISTER_MAX_RETRY_DURATION_SECONDS", 2*time.Minute)
	startTime := time.Now()

	for {
		if time.Since(startTime) > maxRetryDuration {
			logger.Fatal("Failed to register service after %s: %s", maxRetryDuration, serviceName)
			panic("Error registering service after max retries")
		}

		err := consulClient.Agent().ServiceRegister(registration)
		if err != nil {
			logger.Error("Error registering service (elapsed time: %s): %s", time.Since(startTime), err.Error())
			time.Sleep(retryDelay)
		} else {
			logger.Info("Service registered successfully -> {}:{}", serviceName, serviceID)
			return
		}
	}
}

func UnregisterInstance() error {
	if consulClient == nil {
		return fmt.Errorf("onsul client is not initialized")
	}

	logger.Info("Unregistering service -> {}, {}", serviceName, serviceID)

	err := consulClient.Agent().ServiceDeregister(serviceID)
	if err != nil {
		logger.Error("Error unregistering service: {}", err.Error())
		return err
	}

	logger.Info("Service successfully unregistered from Consul")
	return nil
}

func initHeartbeat(port int, host string) {
	go func() {
		duration, err := strconv.ParseInt(utils.GetEnvVariable("DISCOVERY_CLIENT_HEARTBEAT_FREQUENCY_DURATION", "30"), 10, 64)
		if err != nil || duration < 30 {
			duration = 30
		}

		heartbeatFrequency := time.Duration(duration * int64(time.Second))
		time.Sleep(heartbeatFrequency)

		for {
			sendHeartbeat(port, host)
			time.Sleep(heartbeatFrequency)
		}
	}()
}

func sendHeartbeat(port int, host string) {
	logger.Debug("Sending heartbeat -> {}, {}", serviceName, serviceID)

	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		logger.Error("Failed to fetch service status from Consul: {}", err.Error())
		return
	}

	isRegistered := false
	for _, service := range services {
		if service.Service.ID == serviceID {
			isRegistered = true
			break
		}
	}

	if !isRegistered {
		logger.Debug("Service not found in Consul, re-registering...")
		registerService(port, host)
	} else {
		logger.Trace("Service heartbeat successful -> {}, {}", serviceName, serviceID)
	}
}
