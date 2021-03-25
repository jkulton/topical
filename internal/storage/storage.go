package storage

import (
	"bytes"
	"database/sql"
	"github.com/jkulton/topical/internal/models"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"log"
	"time"
)

// Storage is an interface for interacting with a storage layer
type Storage struct {
	db *sql.DB
}

type TopicalStore interface {
	GetTopic(id int) (*models.Topic, error)
	GetRecentTopics() ([]models.Topic, error)
	CreateMessage(m *models.Message) (*models.Message, error)
	CreateTopic(title string) (*models.Topic, error)
}

// New returns a new TopicalStore
func New(db *sql.DB) *Storage {
	return &Storage{db}
}

// GetTopic retrieves a topic from DB by topic
func (t *Storage) GetTopic(id int) (*models.Topic, error) {
	topic := models.Topic{}
	messages := []models.Message{}
	query := `
		SELECT topics.id, topics.title, messages.content, messages.author_initials, messages.author_theme, messages.posted, messages.id
		FROM topics
		INNER JOIN messages ON messages.topic_id = topics.id
		WHERE topics.id = $1
		ORDER BY posted ASC;`

	rows, err := t.db.Query(query, id)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var topicID, authorTheme, messageID int
		var title, content, authorInitials string
		var posted time.Time
		var unsafeHTML bytes.Buffer

		if err = rows.Scan(&topicID, &title, &content, &authorInitials, &authorTheme, &posted, &messageID); err != nil {
			log.Fatal(err)
			return nil, err
		}

		topic.ID = &topicID
		topic.Title = title

		if err := goldmark.Convert([]byte(content), &unsafeHTML); err != nil {
			panic(err)
		}
		safeHTML := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML.Bytes())

		messages = append(messages, models.Message{
			ID:             &messageID,
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
func (t *Storage) GetRecentTopics() ([]models.Topic, error) {
	topics := []models.Topic{}
	query := `
		SELECT DISTINCT topics.*,
			(SELECT COUNT(messages.id) FROM messages WHERE topic_id = topics.id) AS "message_count",
			(SELECT author_initials FROM messages WHERE topic_id = topics.id ORDER BY posted ASC LIMIT 1) AS "author_initials",
			(SELECT author_theme FROM messages WHERE topic_id = topics.id ORDER BY posted ASC LIMIT 1) AS "author_theme",
			(SELECT posted FROM messages WHERE topic_id = topics.id ORDER BY posted DESC LIMIT 1) AS "last_message"
		FROM topics
		INNER JOIN messages
		ON topics.id = messages.topic_id
		ORDER BY last_message DESC
		LIMIT 50;`
	rows, err := t.db.Query(query)

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

		topics = append(topics, models.Topic{
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
func (t *Storage) CreateMessage(m *models.Message) (*models.Message, error) {
	sql := `INSERT INTO messages (topic_id, content, author_initials, author_theme) VALUES ($1, $2, $3, $4)`
	_, err := t.db.Exec(sql, &m.TopicID, m.Content, m.AuthorInitials, m.AuthorTheme)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return m, nil
}

// CreateTopic inserts a new topic into the DB
func (s *Storage) CreateTopic(title string) (*models.Topic, error) {
	id := 0
	err := s.db.QueryRow(`INSERT INTO topics (title) VALUES ($1) RETURNING id`, title).Scan(&id)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return &models.Topic{ID: &id, Title: title}, nil
}
