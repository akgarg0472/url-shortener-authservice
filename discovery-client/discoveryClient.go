package consul

import (
	"fmt"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
)

var logger = Logger.GetLogger("discoveryClient.go")

func InitDiscoveryClient(port int) {
	client := eureka.NewClient([]string{"http://localhost:8761/eureka/v2"})

	instanceInfo := eureka.NewInstanceInfo("localhost", "urlshortener-auth-service", "127.0.0.1", 8081, 60, false)
	instanceInfo.InstanceID = "localhost:urlshortener-auth-service"
	instanceInfo.VipAddress = "urlshortener-auth-service"
	instanceInfo.SecureVipAddress = "urlshortener-auth-service"

	registerInstance(client, instanceInfo)

	initHeartbeat(client, instanceInfo, 30*time.Second)
}

func initHeartbeat(client *eureka.Client, instanceInfo *eureka.InstanceInfo, heartbeatFrequency time.Duration) {
	go func() {
		// Wait for heartbeatFrequency time before sending first heartbeat
		time.Sleep(heartbeatFrequency)

		for {
			logger.Debug("Sending heartbeat to discovery server {}:{}", instanceInfo.App, instanceInfo.InstanceID)
			err := client.SendHeartbeat(instanceInfo.App, instanceInfo.InstanceID)

			if err != nil {
				err, isEurekeError := err.(*eureka.EurekaError)

				if isEurekeError {
					if isInstanceNotFoundError(err) {
						registerInstance(client, instanceInfo)
					}
				}
			}

			time.Sleep(heartbeatFrequency)
		}
	}()
}

func registerInstance(client *eureka.Client, instanceInfo *eureka.InstanceInfo) {
	logger.Debug("Registering instance {}:{}", instanceInfo.App, instanceInfo.InstanceID)

	err := client.RegisterInstance(instanceInfo.App, instanceInfo)

	if err != nil {
		fmt.Println(err)
		return
	}

	logger.Debug("Registered instance with eureka {}:{}", instanceInfo.App, instanceInfo.InstanceID)
}

func isInstanceNotFoundError(err *eureka.EurekaError) bool {
	return err.ErrorCode == 502 && err.Message == "Instance resource not found"
}
