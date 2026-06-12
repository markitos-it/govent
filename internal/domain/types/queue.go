package types

import (
	"errors"
	"log"
	"strings"
	"time"
)

type MessageStatus string

const (
	StatusPending   MessageStatus = "pending"
	StatusProcessed MessageStatus = "processed"
	StatusFailed    MessageStatus = "failed"
)

type QueueMessage struct {
	Id             string        `json:"id" binding:"required,uuid"`
	SubscriberName string        `json:"subscriber_name" binding:"required"`
	EventId        string        `json:"event_id" binding:"required,uuid"`
	Status         MessageStatus `json:"status" binding:"required"`
	CreatedAt      time.Time     `json:"created_at" binding:"required,datetime"`
	UpdatedAt      time.Time     `json:"updated_at" binding:"required,datetime"`
}

func NewQueueMessage(id, subscriberName, eventId string) (*QueueMessage, error) {
	if strings.TrimSpace(id) == "" {
		log.Println("❌ DEBUG ERROR (NewQueueMessage): queue message id cannot be empty")
		return nil, errors.New("queue message id cannot be empty")
	}

	if strings.TrimSpace(subscriberName) == "" {
		log.Println("❌ DEBUG ERROR (NewQueueMessage): subscriber name cannot be empty")
		return nil, errors.New("subscriber name cannot be empty")
	}

	if strings.TrimSpace(eventId) == "" {
		log.Println("❌ DEBUG ERROR (NewQueueMessage): event id cannot be empty")
		return nil, errors.New("event id cannot be empty")
	}

	return &QueueMessage{
		Id:             strings.TrimSpace(id),
		SubscriberName: strings.TrimSpace(subscriberName),
		EventId:        strings.TrimSpace(eventId),
		Status:         StatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (q *QueueMessage) MarkAsProcessed() {
	q.Status = StatusProcessed
	q.UpdatedAt = time.Now()
}

func (q *QueueMessage) MarkAsFailed() {
	q.Status = StatusFailed
	q.UpdatedAt = time.Now()
}

func (QueueMessage) TableName() string {
	return "queue"
}
