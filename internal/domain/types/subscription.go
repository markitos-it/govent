package types

import (
	"errors"
	"log"
	"strings"
)

type Subscription struct {
	Id             string `json:"id" binding:"required,uuid"`
	SubscriberName string `json:"subscriber_name" binding:"required"`
	EventName      string `json:"event_name" binding:"required"`
	Source         string `json:"source" binding:"required"`
}

func NewSubscription(id, subscriberName, eventName, source string) (*Subscription, error) {
	secureId, err := NewSharedId(id)
	if err != nil {
		log.Printf("❌ DEBUG ERROR (NewId): %v\n", err)
		return nil, err
	}

	if strings.TrimSpace(subscriberName) == "" {
		log.Println("❌ DEBUG ERROR (NewSubscription): subscriber name cannot be empty")
		return nil, errors.New("subscriber name cannot be empty")
	}

	if strings.TrimSpace(eventName) == "" {
		log.Println("❌ DEBUG ERROR (NewSubscription): event name cannot be empty")
		return nil, errors.New("event name cannot be empty")
	}

	if strings.TrimSpace(source) == "" {
		log.Println("❌ DEBUG ERROR (NewSubscription): source cannot be empty")
		return nil, errors.New("source cannot be empty")
	}

	return &Subscription{
		Id:             secureId.Value(),
		SubscriberName: strings.TrimSpace(subscriberName),
		EventName:      strings.TrimSpace(eventName),
		Source:         strings.TrimSpace(source),
	}, nil
}

func (Subscription) TableName() string {
	return "subscriptions"
}
