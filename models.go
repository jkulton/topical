package main

// Message represents a message entity, messages have a M:1 relationship with Topics
type Message struct {
	TopicID        *int
	Content        string
	AuthorInitials string
	Posted         string
	AuthorTheme    int
}

// Topic represents a topic entity, one Topic can contain many Messages
type Topic struct {
	ID             *int
	Title          string
	Messages       *[]Message
	MessageCount   *int
	AuthorInitials *string
	AuthorTheme    *string
}

// User is a struct representing a user account, the user stored
// in a simple cookie and defines a name and theme for messages.
type User struct {
	Initials string
	Theme    int
}
