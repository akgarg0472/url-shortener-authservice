package kafka_service

import (
	"context"

	model "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	kafka "github.com/segmentio/kafka-go"
)

var (
	logger                 = Logger.GetLogger("kafkaService.go")
	instance               *KafkaService
	kafkaWriter            *kafka.Writer
	emailNotificationTopic = ""
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
	}
}

func InitKafka() {
	kafkaURL := utils.GetEnvVariable("KAFKA_CONNECTION_URL", "localhost:9092")
	kafkaTopic := getEmailTopic()

	logger.Debug("Initializing Kafka with url: {} and topic: {}", kafkaURL, kafkaTopic)

	kafkaWriter = getKafkaWriter(kafkaURL, kafkaTopic)

	logger.Info("Kafka initialzed: clusterIP={}, topic={}", kafkaWriter.Addr, kafkaWriter.Topic)
}

func CloseKafka() error {
	if kafkaWriter != nil {
		logger.Debug("Closing kafka connection")
		kafkaCloseError := kafkaWriter.Close()
		logger.Info("Kafka connection close status: {}", kafkaCloseError == nil)
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

func (kafkaService *KafkaService) PushNotificationEvent(event model.NotificationEvent) bool {
	logger.Debug("Pushing Event To Kafka topic [{}]: {}", emailNotificationTopic, event.String())

	if kafkaWriter != nil {
		msgBytes, msgError := utils.ConvertToJsonBytes(event)

		if msgError != nil {
			logger.Error("Error converting notification event to bytes: {}", msgError.Error())
			return false
		}

		message := kafka.Message{
			Value: msgBytes,
		}

		msgWriteErr := kafkaWriter.WriteMessages(context.Background(), message)

		return msgWriteErr == nil
	}

	return false
}
