package models

// User is a struct representing a user account, the user stored
// in a simple cookie and defines a name and theme for messages.
type User struct {
	Initials string
	Theme    int
}
