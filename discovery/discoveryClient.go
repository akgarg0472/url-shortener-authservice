package discovery

import (
	"fmt"
	"strconv"
	"time"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

const (
	serviceIDKey = "service_id"
)

var (
	consulClient *api.Client
	serviceID    string
	hostIp       string
	instancePort int
)

func InitDiscoveryClient(port int) {
	isDiscoveryClientEnabled, err := strconv.ParseBool(utils.GetEnvVariable("ENABLE_DISCOVERY_CLIENT", "false"))
	if err != nil || !isDiscoveryClientEnabled {
		if logger.IsInfoEnabled() {
			logger.Info("Discovery client is disabled in configuration")
		}
		return
	}

	consulAddress := utils.GetEnvVariable("DISCOVERY_SERVER_IP", "http://127.0.0.1:8500")

	config := api.DefaultConfig()
	config.Address = consulAddress
	consulClient, err = api.NewClient(config)

	if err != nil {
		if logger.IsFatalEnabled() {
			logger.Fatal("Failed to create Consul client", zap.Error(err))
		}
	}

	hostIp = utils.GetHostIP()
	instancePort = port

	serviceID = fmt.Sprintf("%s-%s", constants.ServiceName, uuid.New().String())

	registerService(false)
	initHeartbeat()
}

func UnregisterInstance() error {
	if consulClient == nil {
		return fmt.Errorf("consul client is not initialized")
	}

	if logger.IsInfoEnabled() {
		logger.Info("Unregistering service",
			zap.String(serviceIDKey, serviceID),
		)
	}

	err := consulClient.Agent().ServiceDeregister(serviceID)

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error("Error unregistering service", zap.Error(err))
		}
		return err
	}

	if logger.IsInfoEnabled() {
		logger.Info("Service successfully unregistered from Consul")
	}

	return nil
}

func registerService(isReRegister bool) {
	if logger.IsInfoEnabled() {
		if isReRegister {
			logger.Info("Re-Registering service with Consul",
				zap.String("service_id", serviceID),
			)
		} else {
			logger.Info("Registering service with Consul",
				zap.String("service_id", serviceID),
			)
		}
	}

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    constants.ServiceName,
		Port:    instancePort,
		Address: hostIp,
		Check: &api.AgentServiceCheck{
			TTL:                            "30s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	retryDelay := utils.GetEnvDurationSeconds("REGISTER_RETRY_DELAY_SECONDS", 5*time.Second)
	maxRetryDuration := utils.GetEnvDurationSeconds("REGISTER_MAX_RETRY_DURATION_SECONDS", 2*time.Minute)
	startTime := time.Now()

	for {
		if time.Since(startTime) > maxRetryDuration {
			if logger.IsFatalEnabled() {
				logger.Fatal("Failed to register service after max retry duration",
					zap.Duration("duration", maxRetryDuration),
				)
			}
			panic("Error registering service after max retries")
		}

		err := consulClient.Agent().ServiceRegister(registration)
		if err != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Error registering service",
					zap.Duration("elapsed_time", time.Since(startTime)),
					zap.Error(err),
				)
			}
			time.Sleep(retryDelay)
		} else {
			if logger.IsInfoEnabled() {
				logger.Info("Service registered successfully",
					zap.String(serviceIDKey, serviceID),
				)
			}
			return
		}
	}
}

func initHeartbeat() {
	go func() {
		duration, err := strconv.ParseInt(utils.GetEnvVariable("DISCOVERY_CLIENT_HEARTBEAT_FREQUENCY_DURATION", "15"), 10, 64)

		if err != nil || duration < 15 {
			duration = 15
		}

		heartbeatFrequency := time.Duration(duration * int64(time.Second))

		for {
			sendHeartbeat()
			time.Sleep(heartbeatFrequency)
		}
	}()
}

func sendHeartbeat() {
	checkID := "service:" + serviceID

	err := consulClient.Agent().UpdateTTL(checkID, "heartbeat passed", api.HealthPassing)

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error("Failed to send heartbeat",
				zap.String(serviceIDKey, serviceID),
				zap.Error(err),
			)
		}

		if statusErr, ok := err.(api.StatusError); ok {
			if statusErr.Code == 404 {
				registerService(true)
			}
		}
	} else {
		if logger.IsDebugEnabled() {
			logger.Debug("Heartbeat updated",
				zap.String(serviceIDKey, serviceID),
			)
		}
	}
}
