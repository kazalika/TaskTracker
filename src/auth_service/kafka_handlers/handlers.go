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

const (
	// It should be some username that will never be used by users or for some another reason
	accountForCreatingEmptyStatistics = "ACCOUNT_FOR_CREATING_EMPTY_STATISTICS"
	// It's message which kafka send sometimes. When it occures should retry call
	kafkaLeadershipErrorMessage = "[5] Leader Not Available: the cluster is in the middle of a leadership election and there is currently no leader for this partition and hence it is unavailable for writes"
)

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}
}

func InitKafkaTopics() {
	kafkaURL := os.Getenv("KAFKA_URL")
	log.Printf("Kafka's URL = %v", kafkaURL)

	views = getKafkaWriter(kafkaURL, "views")
	likes = getKafkaWriter(kafkaURL, "likes")
}

func CreateEmptyStatistics(taskID int32, taskAuthor string) error {
	if err := Like(accountForCreatingEmptyStatistics, taskID, taskAuthor); err != nil {
		return err
	}
	return View(accountForCreatingEmptyStatistics, taskID, taskAuthor)
}

func Like(liker string, taskID int32, taskAuthor string) error {
	encoded, err := json.Marshal(map[string]any{
		"username":    liker,
		"task_id":     taskID,
		"task_author": taskAuthor, // for statistics
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
		if err.Error() != kafkaLeadershipErrorMessage {
			return err
		}
	}
}

func View(viewer string, taskID int32, taskAuthor string) error {
	encoded, err := json.Marshal(map[string]any{
		"username":    viewer,
		"task_id":     taskID,
		"task_author": taskAuthor, // for statistics
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
		if err.Error() != kafkaLeadershipErrorMessage {
			return err
		}
	}
}

func CloseKafkaTopics() {
	views.Close()
	likes.Close()
}
