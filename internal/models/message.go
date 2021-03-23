package models

import "time"

// Message represents a message entity, messages have a M:1 relationship with Topics
type Message struct {
	ID             *int
	TopicID        *int
	Content        string
	AuthorInitials string
	Posted         time.Time
	AuthorTheme    int
}
