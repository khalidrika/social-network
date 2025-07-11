package models

import "time"

type Event struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"groupId"`
	AuthorID    string    `json:"authorId"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Date        string    `json:"date"`        // Keep as string for compatibility
	DateTime    time.Time `json:"dateTime"`    // For proper time handling
	CreatedAt   time.Time `json:"createdAt"`
	
	// Author info for display
	Author User `json:"author"`
}

type EventResponse struct {
	EventID  string `json:"eventId"`
	UserID   string `json:"userId"`
	Response string `json:"response"` // "going", "not_going", "maybe"
	UserName string `json:"userName"`
}

type EventWithResponses struct {
	Event
	Responses    []EventResponse `json:"responses"`
	UserResponse string          `json:"userResponse"` // current user's response
	GoingCount   int            `json:"goingCount"`
	NotGoingCount int           `json:"notGoingCount"`
	MaybeCount   int            `json:"maybeCount"`
}

type EventRepository interface {
	GetAll(groupId string) ([]Event, error)
	GetData(eventID string) (Event, error)
	Save(Event) error
	AddParticipant(eventID, userID string) error
	RemoveParticipant(eventID, userID string) error
	IsParticipating(eventID, userID string) (bool, error)
	
	// New methods for RSVP functionality
	UpdateResponse(eventID, userID, response string) error
	GetEventWithResponses(eventID, currentUserID string) (*EventWithResponses, error)
	GetGroupEventsWithResponses(groupID, currentUserID string) ([]EventWithResponses, error)
}
