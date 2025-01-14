package kafka_service

import (
	"context"

	"github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"github.com/segmentio/kafka-go"
)

var (
	logger                    = Logger.GetLogger("kafkaService.go")
	instance                  *KafkaService
	emailKafkaWriter          *kafka.Writer
	userRegisteredKafkaWriter *kafka.Writer
	emailNotificationTopic    = ""
	userRegisteredTopic       = ""
)

type KafkaService struct {
	emailTopic string
}

func GetInstance() *KafkaService {
	if instance == nil {
		instance = &KafkaService{
			emailTopic: getEmailTopic(),
		}
	}

	return instance
}

func getKafkaWriter(kafkaURL string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		Async:                  true,
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireOne,
	}
}

func InitKafka() {
	kafkaURL := utils.GetEnvVariable("KAFKA_CONNECTION_URL", "localhost:9092")
	emailKafkaTopic := getEmailTopic()
	userRegisteredKafkaTopic := getUserRegisteredTopic()

	logger.Info("Initializing Kafka with url: {} and topic(s): {}", kafkaURL, []string{emailKafkaTopic, userRegisteredKafkaTopic})

	emailKafkaWriter = getKafkaWriter(kafkaURL, emailKafkaTopic)
	userRegisteredKafkaWriter = getKafkaWriter(kafkaURL, userRegisteredKafkaTopic)

	logger.Info("Kafka initialized (email): clusterIP={}, topic={}", emailKafkaWriter.Addr, emailKafkaWriter.Topic)
	logger.Info("Kafka initialized (userRegistered): clusterIP={}, topic={}", userRegisteredKafkaWriter.Addr, userRegisteredKafkaWriter.Topic)
}

func CloseKafka() error {
	if emailKafkaWriter != nil {
		logger.Debug("Closing email kafka connection")
		kafkaCloseError := emailKafkaWriter.Close()
		logger.Info("Email Kafka connection close status: {}", kafkaCloseError == nil)
		return kafkaCloseError
	}

	return nil
}

func getEmailTopic() string {
	emailNotificationTopic = utils.GetEnvVariable("KAFKA_TOPIC_EMAIL_NOTIFICATION", "")

	if emailNotificationTopic == "" {
		panic("KAFKA_TOPIC_EMAIL_NOTIFICATION not found")
	}

	return emailNotificationTopic
}

func getUserRegisteredTopic() string {
	userRegisteredTopic = utils.GetEnvVariable("KAFKA_TOPIC_USER_REGISTERED", "")

	if userRegisteredTopic == "" {
		panic("KAFKA_TOPIC_USER_REGISTERED not found")
	}

	return userRegisteredTopic
}

func (kafkaService *KafkaService) PushNotificationEvent(reqId string, event model.NotificationEvent) {
	logger.Debug("[{}] Pushing Notification Event To Kafka topic [{}]: {}", reqId, emailNotificationTopic, event.String())

	if emailKafkaWriter != nil {
		msgBytes, msgError := utils.ConvertToJsonBytes(event)

		if msgError != nil {
			logger.Error("[{}] Error converting notification event to bytes: {}", reqId, msgError.Error())
			return
		}

		message := kafka.Message{
			Value: msgBytes,
		}

		msgWriteErr := emailKafkaWriter.WriteMessages(context.Background(), message)

		logger.Debug("[{}] Kafka push result: {}", reqId, msgWriteErr == nil)
	}
}

func (kafkaService *KafkaService) PushUserRegisteredEvent(reqId string, userId string) {
	logger.Debug("[{}] Pushing user Registered Event To Kafka topic [{}]: {}", reqId, userRegisteredTopic, userId)

	if userRegisteredKafkaWriter != nil {
		event := model.UserRegisteredEvent{
			UserId: userId,
		}

		msgBytes, msgError := utils.ConvertToJsonBytes(event)

		if msgError != nil {
			logger.Error("[{}] Error converting user registered event to bytes: {}", reqId, msgError.Error())
			return
		}

		message := kafka.Message{
			Value: msgBytes,
		}

		msgWriteErr := userRegisteredKafkaWriter.WriteMessages(context.Background(), message)

		logger.Debug("[{}] Kafka push result: {}", reqId, msgWriteErr == nil)

	}
}
