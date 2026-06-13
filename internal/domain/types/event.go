package types

import (
	"log"
	"time"
)

type Event struct {
	Id        string    `json:"id" binding:"required,uuid"`
	Name      string    `json:"name" binding:"required"`
	Source    string    `json:"source" binding:"required"`
	Payload   string    `json:"payload" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required,datetime" default:"now"`
	UpdatedAt time.Time `json:"updated_at" binding:"required,datetime" default:"now"`
}

func NewEvent(id, name, source, payload string) (*Event, error) {
	secureId, err := NewSharedId(id)
	if err != nil {
		log.Printf("❌ DEBUG ERROR (Id): %v\n", err)
		return nil, err
	}

	secureName, err := NewSlug(name)
	if err != nil {
		log.Printf("❌ DEBUG ERROR (Slug): %v\n", err)
		return nil, err
	}

	secureSource, err := NewSource(source)
	if err != nil {
		log.Printf("❌ DEBUG ERROR (Source): %v\n", err)
		return nil, err
	}

	securePayload, err := NewPayload(payload)
	if err != nil {
		log.Printf("❌ DEBUG ERROR (Payload): %v\n", err)
		return nil, err
	}

	return &Event{
		Id:        secureId.Value(),
		Name:      secureName.Value(),
		Source:    secureSource.Value(),
		Payload:   securePayload.Value(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (e *Event) GetName() *Name {
	result, _ := NewName(e.Name)

	return result
}

func (e *Event) GetSource() *Source {
	result, _ := NewSource(e.Source)

	return result
}

func (Event) TableName() string {
	return "events"
}
