package consul

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	logger          = Logger.GetLogger("discoveryClient.go")
	discoveryClient *eureka.Client
	instanceInfo    *eureka.InstanceInfo
)

func InitDiscoveryClient(port int) {
	isDiscoveryClientEnabled, err := strconv.ParseBool(utils.GetEnvVariable("ENABLE_DISCOVERY_CLIENT", "false"))

	if err != nil || !isDiscoveryClientEnabled {
		logger.Info("Discovery client is disabled in configuration")
		return
	}

	discoveryServerIp := strings.Split(utils.GetEnvVariable("DISCOVERY_SERVER_IP", "http://localhost:8761/eureka/v2"), ",")

	discoveryClient = eureka.NewClient(discoveryServerIp)

	appId := "urlshortener-auth-service"
	host := utils.GetHostIP()

	instanceId := fmt.Sprintf("%s:%s:%d", generateRandomInstanceId(), appId, port)
	appAddress := "urlshortener-auth-service"

	instanceInfo = eureka.NewInstanceInfo(host, appId, host, port, 60, false)
	instanceInfo.InstanceID = instanceId
	instanceInfo.VipAddress = appAddress
	instanceInfo.SecureVipAddress = appAddress
	instanceInfo.HealthCheckUrl = fmt.Sprintf("http://%s:%d/admin/health", host, port)
	instanceInfo.StatusPageUrl = fmt.Sprintf("http://%s:%d/admin/info", host, port)

	registerInstance()
	initHeartbeat()
}

func UnregisterInstance() error {
	if discoveryClient == nil {
		return fmt.Errorf("discovery client is not initialized")
	}

	logger.Info("unregistering instance -> {}, {}", instanceInfo.App, instanceInfo.InstanceID)

	deleteEndpoint := fmt.Sprintf("apps/%s/%s", instanceInfo.App, instanceInfo.InstanceID)

	_, err := discoveryClient.Delete(deleteEndpoint)

	if err != nil {
		logger.Error("error unregistering instance: {}", err.Error())
		return err
	}

	return nil
}

func initHeartbeat() {
	go func() {
		duration, err := strconv.ParseInt(utils.GetEnvVariable("DISCOVERY_CLIENT_HEARTBEAT_FREQUENCY_DURATION", "30"), 10, 64)

		if err != nil || duration < 30 {
			duration = 30
		}

		heartbeatFrequency := time.Duration(duration * int64(time.Second))

		time.Sleep(heartbeatFrequency)

		for {
			sendHeartbeat()
			time.Sleep(heartbeatFrequency)
		}
	}()
}

func sendHeartbeat() {
	logger.Debug("sending heartbeat -> {}, {}", instanceInfo.App, instanceInfo.InstanceID)
	err := discoveryClient.SendHeartbeat(instanceInfo.App, instanceInfo.InstanceID)

	if err != nil {
		var err *eureka.EurekaError
		isEurekaError := errors.As(err, &err)

		if isEurekaError {
			if isInstanceNotFoundError(err) {
				registerInstance()
			}
		}
	} else {
		logger.Trace("heartbeat sent successfully -> {}, {}", instanceInfo.App, instanceInfo.InstanceID)
	}
}

func registerInstance() {
	go func() {
		logger.Info("registering instance -> {}:{}", instanceInfo.App, instanceInfo.InstanceID)

		retryDelay := utils.GetEnvDurationSeconds("REGISTER_RETRY_DELAY_SECONDS", 5*time.Second)
		maxRetryDuration := utils.GetEnvDurationSeconds("REGISTER_MAX_RETRY_DURATION_SECONDS", 2*time.Minute)

		var startTime = time.Now()

		for {
			elapsed := time.Since(startTime)

			if elapsed > maxRetryDuration {
				logger.Fatal("Failed to register instance after %s: %s:%s", maxRetryDuration, instanceInfo.App, instanceInfo.InstanceID)
				panic("Error registering instance after max retries")
			}

			err := discoveryClient.RegisterInstance(instanceInfo.App, instanceInfo)

			if err != nil {
				logger.Error("Error registering instance (elapsed time: %s): %s", elapsed, err.Error())
				time.Sleep(retryDelay)
			} else {
				logger.Info("Instance registered successfully -> {}:{}", instanceInfo.App, instanceInfo.InstanceID)
				return
			}
		}
	}()
}

func isInstanceNotFoundError(err *eureka.EurekaError) bool {
	return err != nil && err.ErrorCode == 502 && err.Message == "Instance resource not found"
}

func generateRandomInstanceId() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 12
	var result strings.Builder
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return utils.GetHostIP()
	}

	for _, b := range randomBytes {
		result.WriteByte(charset[b%byte(len(charset))])
	}

	return result.String()
}
