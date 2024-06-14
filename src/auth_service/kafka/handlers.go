package kafka_handlers

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

var (
	views *kafka.Writer
	likes *kafka.Writer
)

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}
}

func InitKafkaConnections() {
	kafkaURL := os.Getenv("KAFKA_URL")
	log.Printf("Kafka's URL = %v", kafkaURL)

	views = getKafkaWriter(kafkaURL, "views")
	likes = getKafkaWriter(kafkaURL, "likes")
}

func CreateEmptyStatistics(taskID string, taskAuthor string) error {
	if err := Like("ADMIN_NICKNAME_HERE", taskID, taskAuthor); err != nil {
		return err
	}
	return View("ADMIN_NICKNAME_HERE", taskID, taskAuthor)
}

func Like(liker string, taskID string, taskAuthor string) error {
	encoded, err := json.Marshal(map[string]string{
		"username":    liker,
		"task_id":     taskID,
		"task_author": taskAuthor,
	})
	if err != nil {
		return err
	}
	requestID := uuid.New().String()

	log.Printf("Send message (like) to Kafka {Key: %s, Value: %s}", requestID, string(encoded))
	for {
		err = likes.WriteMessages(context.Background(), kafka.Message{Key: []byte(requestID), Value: encoded})
		if err == nil {
			return nil
		}
		if err.Error() != "[5] Leader Not Available: the cluster is in the middle of a leadership election and there is currently no leader for this partition and hence it is unavailable for writes" {
			return err
		}
	}
}

func View(viewer string, taskID string, taskAuthor string) error {
	encoded, err := json.Marshal(map[string]string{
		"username":    viewer,
		"task_id":     taskID,
		"task_author": taskAuthor,
	})
	if err != nil {
		return err
	}
	requestID := uuid.New().String()

	log.Printf("Send message (view) to Kafka {Key: %s, Value: %s}", requestID, string(encoded))
	for {
		err = views.WriteMessages(context.Background(), kafka.Message{Key: []byte(requestID), Value: encoded})
		if err == nil {
			return nil
		}
		if err.Error() != "[5] Leader Not Available: the cluster is in the middle of a leadership election and there is currently no leader for this partition and hence it is unavailable for writes" {
			return err
		}
	}
}

func CloseKafkaConnections() {
	views.Close()
	likes.Close()
}