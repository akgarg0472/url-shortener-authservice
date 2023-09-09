package consul

import (
	"fmt"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
)

var (
	logger          = Logger.GetLogger("discoveryClient.go")
	discoveryClient *eureka.Client
	instanceInfo    *eureka.InstanceInfo
)

func InitDiscoveryClient(port int) {
	discoveryClient = eureka.NewClient([]string{"http://localhost:8761/eureka/v2"})

	host := "localhost"
	appId := "urlshortener-auth-service"
	instanceId := fmt.Sprintf("%s:urlshortener-auth-service:%d", host, port)
	appAddress := "urlshortener-auth-service"

	instanceInfo = eureka.NewInstanceInfo(host, appId, "127.0.0.1", port, 60, false)
	instanceInfo.InstanceID = instanceId
	instanceInfo.VipAddress = appAddress
	instanceInfo.SecureVipAddress = appAddress

	err := registerInstance()

	if err == nil {
		initHeartbeat(30 * time.Second)
	}
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

func initHeartbeat(heartbeatFrequency time.Duration) {
	go func() {
		// Wait for heartbeatFrequency time before sending first heartbeat
		time.Sleep(heartbeatFrequency)

		for {
			logger.Debug("sending heartbeat -> {}, {}", instanceInfo.App, instanceInfo.InstanceID)
			err := discoveryClient.SendHeartbeat(instanceInfo.App, instanceInfo.InstanceID)

			if err != nil {
				err, isEurekeError := err.(*eureka.EurekaError)

				if isEurekeError {
					if isInstanceNotFoundError(err) {
						registerInstance()
					}
				}
			}

			time.Sleep(heartbeatFrequency)
		}
	}()
}

func registerInstance() error {
	logger.Info("registering instance -> {}:{}", instanceInfo.App, instanceInfo.InstanceID)

	err := discoveryClient.RegisterInstance(instanceInfo.App, instanceInfo)

	if err != nil {
		logger.Error("error registering instance: {}", err.Error())
		return err
	}

	logger.Debug("instance registered -> {}:{}", instanceInfo.App, instanceInfo.InstanceID)

	return nil
}

func isInstanceNotFoundError(err *eureka.EurekaError) bool {
	return err.ErrorCode == 502 && err.Message == "Instance resource not found"
}
