package kafka_service

import (
	"context"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var (
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

	if logger.IsInfoEnabled() {
		logger.Info("Initializing Kafka with url and topic(s)",
			zap.String("kafka_url", kafkaURL),
			zap.Strings("topics", []string{emailKafkaTopic, userRegisteredKafkaTopic}),
		)
	}

	emailKafkaWriter = getKafkaWriter(kafkaURL, emailKafkaTopic)
	userRegisteredKafkaWriter = getKafkaWriter(kafkaURL, userRegisteredKafkaTopic)

	if logger.IsInfoEnabled() {
		logger.Info("Kafka initialized (email)",
			zap.Any("clusterIP", emailKafkaWriter.Addr),
			zap.String("topic", emailKafkaWriter.Topic),
		)
	}
	if logger.IsInfoEnabled() {
		logger.Info("Kafka initialized (userRegistered)",
			zap.Any("clusterIP", userRegisteredKafkaWriter.Addr),
			zap.String("topic", userRegisteredKafkaWriter.Topic),
		)
	}
}

func CloseKafka() error {
	if emailKafkaWriter != nil {
		logger.Debug("Closing email kafka connection")
		kafkaCloseError := emailKafkaWriter.Close()
		if logger.IsInfoEnabled() {
			logger.Info("Email Kafka connection closed",
				zap.Bool("status", kafkaCloseError == nil),
			)
		}
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
	if logger.IsDebugEnabled() {
		logger.Debug(
			"Pushing Notification Event To Kafka",
			zap.String(constants.RequestIdLogKey, reqId),
			zap.String("topic", emailNotificationTopic),
			zap.String("event", event.String()),
		)
	}

	if emailKafkaWriter != nil {
		msgBytes, msgError := utils.ConvertToJsonBytes(event)

		if msgError != nil {
			if logger.IsErrorEnabled() {
				logger.Error(
					"Error converting notification event to bytes",
					zap.String(constants.RequestIdLogKey, reqId),
					zap.Error(msgError),
				)
			}
			return
		}

		message := kafka.Message{
			Value: msgBytes,
		}

		msgWriteErr := emailKafkaWriter.WriteMessages(context.Background(), message)

		if logger.IsDebugEnabled() {
			logger.Debug(
				"Kafka push result",
				zap.String(constants.RequestIdLogKey, reqId),
				zap.Bool("success", msgWriteErr == nil),
			)
		}
	}
}

func (kafkaService *KafkaService) PushUserRegisteredEvent(reqId string, userId string) {
	if logger.IsDebugEnabled() {
		logger.Debug(
			"Pushing user Registered Event To Kafka",
			zap.String(constants.RequestIdLogKey, reqId),
			zap.String("topic", userRegisteredTopic),
			zap.String("userId", userId),
		)
	}

	if userRegisteredKafkaWriter != nil {
		event := model.UserRegisteredEvent{
			UserId: userId,
		}

		msgBytes, msgError := utils.ConvertToJsonBytes(event)

		if msgError != nil {
			if logger.IsErrorEnabled() {
				logger.Error(
					"Error converting user registered event to bytes",
					zap.String(constants.RequestIdLogKey, reqId),
					zap.String("error", msgError.Error()),
				)
			}
			return
		}

		message := kafka.Message{
			Value: msgBytes,
		}

		msgWriteErr := userRegisteredKafkaWriter.WriteMessages(context.Background(), message)

		if logger.IsDebugEnabled() {
			logger.Debug(
				"Kafka push result",
				zap.String(constants.RequestIdLogKey, reqId),
				zap.Bool("success", msgWriteErr == nil),
			)
		}
	}
}
