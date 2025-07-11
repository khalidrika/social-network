package db

import (
	"database/sql"
	"time"
	"social-network/pkg/models"
)

type EventRepository struct {
	DB *sql.DB
}

func (repo *EventRepository) GetAll(groupID string) ([]models.Event, error) {
	var events = []models.Event{}
	rows, err := repo.DB.Query(`
		SELECT 
			e.event_id, 
			e.group_id, 
			e.created_by, 
			e.title, 
			e.content, 
			e.date,
			e.created_at,
			u.first_name,
			u.last_name,
			u.email
		FROM event e 
		JOIN users u ON e.created_by = u.user_id 
		WHERE e.group_id = ? 
		ORDER BY e.date DESC`, groupID)
	if err != nil {
		return events, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var event models.Event
		var dateStr string
		var createdAtStr string
		err := rows.Scan(
			&event.ID, 
			&event.GroupID, 
			&event.AuthorID, 
			&event.Title, 
			&event.Content, 
			&dateStr,
			&createdAtStr,
			&event.Author.FirstName,
			&event.Author.LastName,
			&event.Author.Email,
		)
		if err != nil {
			continue
		}
		
		// Parse the date string
		if dateTime, err := time.Parse(time.RFC3339, dateStr); err == nil {
			event.DateTime = dateTime
			event.Date = dateStr // Keep original ISO format for frontend parsing
		}
		
		// Parse created_at
		if createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			event.CreatedAt = createdAt
		}
		
		// Set author ID
		event.Author.ID = event.AuthorID
		
		events = append(events, event)
	}
	return events, nil
}

func (repo *EventRepository) GetData(eventId string) (models.Event, error) {
	row := repo.DB.QueryRow(`
		SELECT 
			e.title, 
			e.content, 
			e.event_id, 
			e.group_id, 
			e.date, 
			e.created_by,
			e.created_at,
			u.first_name,
			u.last_name,
			u.email
		FROM event e 
		JOIN users u ON e.created_by = u.user_id 
		WHERE e.event_id = ?`, eventId)
	
	var event models.Event
	var dateStr string
	var createdAtStr string
	
	err := row.Scan(
		&event.Title, 
		&event.Content, 
		&event.ID, 
		&event.GroupID, 
		&dateStr, 
		&event.AuthorID,
		&createdAtStr,
		&event.Author.FirstName,
		&event.Author.LastName,
		&event.Author.Email,
	)
	if err != nil {
		return event, err
	}
	
	// Parse the date string to preserve time
	if dateTime, err := time.Parse(time.RFC3339, dateStr); err == nil {
		event.DateTime = dateTime
		event.Date = dateStr // Keep original ISO format for frontend parsing
	}
	
	// Parse created_at
	if createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
		event.CreatedAt = createdAt
	}
	
	// Set author ID
	event.Author.ID = event.AuthorID
	
	return event, nil
}

func (repo *EventRepository) Save(event models.Event) error {
	stmt, err := repo.DB.Prepare("INSERT INTO event (event_id, group_id, created_by, content, title, date) values (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	// Use DateTime (time.Time) instead of Date (string) for proper time storage
	if _, err := stmt.Exec(event.ID, event.GroupID, event.AuthorID, event.Content, event.Title, event.DateTime); err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) AddParticipant(eventID, userID string) error {
	stmt, err := repo.DB.Prepare("INSERT INTO event_users (event_id, user_id) values (?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(eventID, userID); err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) RemoveParticipant(eventID, userID string) error {
	stmt, err := repo.DB.Prepare("DELETE FROM event_users WHERE user_id = ? AND event_id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) UpdateResponse(eventID, userID, response string) error {
	// First check if user already has a response
	var count int
	row := repo.DB.QueryRow("SELECT COUNT(*) FROM event_users WHERE event_id = ? AND user_id = ?", eventID, userID)
	if err := row.Scan(&count); err != nil {
		return err
	}
	
	if count > 0 {
		// Update existing response
		stmt, err := repo.DB.Prepare("UPDATE event_users SET response = ? WHERE event_id = ? AND user_id = ?")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(response, eventID, userID)
		return err
	} else {
		// Insert new response
		stmt, err := repo.DB.Prepare("INSERT INTO event_users (event_id, user_id, response) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(eventID, userID, response)
		return err
	}
}

func (repo *EventRepository) GetEventWithResponses(eventID, currentUserID string) (*models.EventWithResponses, error) {
	// Get event details
	event, err := repo.GetData(eventID)
	if err != nil {
		return nil, err
	}
	
	// Get author info
	authorRow := repo.DB.QueryRow("SELECT first_name, last_name FROM users WHERE user_id = ?", event.AuthorID)
	var firstName, lastName string
	if err := authorRow.Scan(&firstName, &lastName); err == nil {
		event.Author.FirstName = firstName
		event.Author.LastName = lastName
	}
	
	eventWithResponses := &models.EventWithResponses{
		Event: event,
		Responses: []models.EventResponse{},
	}
	
	// Get all responses
	rows, err := repo.DB.Query(`
		SELECT eu.user_id, eu.response, u.first_name, u.last_name 
		FROM event_users eu 
		JOIN users u ON eu.user_id = u.user_id 
		WHERE eu.event_id = ?`, eventID)
	if err != nil {
		return eventWithResponses, nil // Return event even if no responses
	}
	defer rows.Close()
	
	for rows.Next() {
		var response models.EventResponse
		var firstName, lastName string
		rows.Scan(&response.UserID, &response.Response, &firstName, &lastName)
		response.EventID = eventID
		response.UserName = firstName + " " + lastName
		
		eventWithResponses.Responses = append(eventWithResponses.Responses, response)
		
		// Set current user's response
		if response.UserID == currentUserID {
			eventWithResponses.UserResponse = response.Response
		}
		
		// Count responses
		switch response.Response {
		case "going":
			eventWithResponses.GoingCount++
		case "not_going":
			eventWithResponses.NotGoingCount++
		case "maybe":
			eventWithResponses.MaybeCount++
		}
	}
	
	return eventWithResponses, nil
}

func (repo *EventRepository) GetGroupEventsWithResponses(groupID, currentUserID string) ([]models.EventWithResponses, error) {
	events, err := repo.GetAll(groupID)
	if err != nil {
		return nil, err
	}
	
	var eventsWithResponses []models.EventWithResponses
	
	for _, event := range events {
		eventWithResp, err := repo.GetEventWithResponses(event.ID, currentUserID)
		if err != nil {
			continue // Skip this event if there's an error
		}
		eventsWithResponses = append(eventsWithResponses, *eventWithResp)
	}
	
	return eventsWithResponses, nil
}

func (repo *EventRepository) IsParticipating(eventID, userID string) (bool, error) {
	row := repo.DB.QueryRow("SELECT COUNT() FROM event_users WHERE event_id = ? AND user_id = ?", eventID, userID)
	var participate int
	if err := row.Scan(&participate); err != nil {
		return false, err
	}
	return participate > 0, nil
}
