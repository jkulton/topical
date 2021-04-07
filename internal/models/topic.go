package models

// Topic represents a topic entity, one Topic can contain many Messages
type Topic struct {
	ID             *int
	Title          string
	Messages       *[]Message
	MessageCount   *int
	AuthorInitials *string
	AuthorTheme    *string
}
