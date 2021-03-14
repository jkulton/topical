package main

import (
	"bytes"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"log"
)

// Storage is an interface for interacting with a storage layer
type Storage struct {
	db *sql.DB
}

// NewStorage creates a new Storage entity with a provided driver and database name
func NewStorage(driver, database string) (*Storage, error) {
	var err error

	stg := new(Storage)
	stg.db, err = sql.Open(driver, database)

	if err != nil {
		return nil, err
	}

	return stg, nil
}

// GetTopic retrieves a topic from DB by topic
func (s *Storage) GetTopic(id int) (*Topic, error) {
	topic := Topic{}
	messages := []Message{}
	query := `
		SELECT
			topics.id,
			topics.title,
			messages.content,
			messages.author_initials,
			messages.author_theme,
			messages.posted
		FROM topics
		INNER JOIN messages ON messages.topic_id = topics.id
		WHERE topics.id = ?
		ORDER BY posted;
	`

	rows, err := s.db.Query(query, id)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		var content string
		var authorInitials string
		var posted string
		var authorTheme int
		var unsafeHTML bytes.Buffer

		if err = rows.Scan(&id, &title, &content, &authorInitials, &authorTheme, &posted); err != nil {
			log.Fatal(err)
			return nil, err
		}

		topic.ID = &id
		topic.Title = title

		if err := goldmark.Convert([]byte(content), &unsafeHTML); err != nil {
			panic(err)
		}
		safeHTML := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML.Bytes())

		messages = append(messages, Message{
			Content:        string(safeHTML),
			AuthorInitials: authorInitials,
			Posted:         posted,
			AuthorTheme:    authorTheme,
		})
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	topic.Messages = &messages

	return &topic, nil
}

// GetRecentTopics returns a list of the 50 most recently posted-on topics
func (s *Storage) GetRecentTopics() ([]Topic, error) {
	topics := []Topic{}
	query := `
		SELECT DISTINCT topics.*,
			(SELECT COUNT(messages.id) FROM messages WHERE topic_id = topics.id) AS "message_count",
			(SELECT author_initials FROM messages WHERE topic_id = topics.id ORDER BY posted LIMIT 1) AS "author_initials",
			(SELECT author_theme FROM messages WHERE topic_id = topics.id ORDER BY posted LIMIT 1) AS "author_theme",
			(SELECT posted FROM messages WHERE topic_id = topics.id ORDER BY posted DESC LIMIT 1) AS "last_message"
		FROM topics
		INNER JOIN messages
		ON topics.id = messages.topic_id
		ORDER BY last_message DESC
		LIMIT 50;
	`
	rows, err := s.db.Query(query)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id, messageCount int
		var title, authorInitials, authorTheme, posted string
		err = rows.Scan(&id, &title, &messageCount, &authorInitials, &authorTheme, &posted)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		topics = append(topics, Topic{
			ID:             &id,
			Title:          title,
			MessageCount:   &messageCount,
			AuthorInitials: &authorInitials,
			AuthorTheme:    &authorTheme,
		})
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return topics, nil
}

// CreateMessage inserts a message into the DB
func (s *Storage) CreateMessage(m *Message) (*Message, error) {
	sql := `INSERT INTO messages (topic_id, content, author_initials, author_theme) VALUES (?, ?, ?, ?)`
	query, err := s.db.Prepare(sql)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	_, err = query.Exec(&m.TopicID, m.Content, m.AuthorInitials, m.AuthorTheme)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return m, nil
}

// CreateTopic inserts a new topic into the DB
func (s *Storage) CreateTopic(title string) (*Topic, error) {
	query, err := s.db.Prepare(`INSERT INTO topics (title) VALUES (?)`)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	res, err := query.Exec(title)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	idInt64, err := res.LastInsertId()

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	id := int(idInt64)

	return &Topic{ID: &id, Title: title}, nil
}
