package consul

import (
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

	discClientMachinesIP := strings.Split(utils.GetEnvVariable("DISCOVERY_CLIENT_IP", "http://localhost:8761/eureka/v2"), ",")

	discoveryClient = eureka.NewClient(discClientMachinesIP)

	host := "localhost"
	appId := "urlshortener-auth-service"
	instanceId := fmt.Sprintf("%s:urlshortener-auth-service:%d", host, port)
	appAddress := "urlshortener-auth-service"

	instanceInfo = eureka.NewInstanceInfo(host, appId, "127.0.0.1", port, 60, false)
	instanceInfo.InstanceID = instanceId
	instanceInfo.VipAddress = appAddress
	instanceInfo.SecureVipAddress = appAddress

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
	logger.Info("registering instance -> {}:{}", instanceInfo.App, instanceInfo.InstanceID)

	err := discoveryClient.RegisterInstance(instanceInfo.App, instanceInfo)

	if err != nil {
		logger.Error("error registering instance: {}", err.Error())
		return
	}

	logger.Debug("instance registered -> {}:{}", instanceInfo.App, instanceInfo.InstanceID)
}

func isInstanceNotFoundError(err *eureka.EurekaError) bool {
	return err.ErrorCode == 502 && err.Message == "Instance resource not found"
}
