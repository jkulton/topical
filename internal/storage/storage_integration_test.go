package storage

import (
	"context"
	"database/sql"
	"github.com/jkulton/topical/internal/models"
	_ "github.com/lib/pq" // Postgres driver
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io/ioutil"
	"strings"
	"testing"
)

var db *sql.DB
var dbContainer testcontainers.Container

type TestHelper struct {
	DB          *sql.DB
	DBContainer testcontainers.Container
	Context     context.Context
}

func testSetup() (th TestHelper) {
	th.Context = context.Background()
	th.DBContainer, _ = createPostgresContainer(th.Context)
	th.DB, _ = createTestDB(th.DBContainer, th.Context)
	return th
}

func testTeardown(th TestHelper) {
	th.DB.Close()
	th.DBContainer.Terminate(th.Context)
}

func createPostgresContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:13.2-alpine",
		Env: map[string]string{
			"POSTGRES_DB":       "test",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	return dbContainer, err
}

func createTestDB(dbContainer testcontainers.Container, ctx context.Context) (*sql.DB, error) {
	endpoint, _ := dbContainer.Endpoint(ctx, "")
	endpoint = "postgresql://test:test@" + endpoint + "?sslmode=disable"
	db, err := sql.Open("postgres", endpoint)

	if err != nil {
		return nil, err
	}

	content, _ := ioutil.ReadFile("../../schema.sql")

	schema := string(content)
	db.Exec(schema)

	content, _ = ioutil.ReadFile("../../seeds.sql")
	seeds := string(content)
	db.Exec(seeds)

	return db, nil
}

func TestCreateMessageIntegration(t *testing.T) {
	t.Run("creates a message in an existing topic", func(t *testing.T) {
		th := testSetup()
		content := "Test Message"
		store := New(th.DB)

		topics, _ := store.GetRecentTopics()
		topicID := topics[0].ID
		message, _ := store.CreateMessage(&models.Message{TopicID: topicID, Content: content, AuthorInitials: "JK", AuthorTheme: 1})
		topic, _ := store.GetTopic(*message.TopicID)
		lastMessageInTopic := (*topic.Messages)[len(*topic.Messages)-1]

		if strings.Contains(lastMessageInTopic.Content, content) == false {
			t.Error("expected new message to be last message in topic")
		}

		testTeardown(th)
	})
}

func TestGetTopicIntegration(t *testing.T) {
	t.Run("returns existing topic", func(t *testing.T) {
		th := testSetup()

		store := New(th.DB)
		topics, _ := store.GetRecentTopics()

		firstTopic := topics[0]
		topicID := firstTopic.ID

		topic, _ := store.GetTopic(*topicID)

		if topic.Title != firstTopic.Title {
			t.Error("expected topic not returned")
		}

		testTeardown(th)
	})
}

func TestGetRecentTopicsIntegration(t *testing.T) {
	t.Run("returns list of recent topics", func(t *testing.T) {
		th := testSetup()

		topics, _ := New(th.DB).GetRecentTopics()

		if len(topics) != 3 {
			t.Error("expected three recent topics")
		}

		testTeardown(th)
	})
}

func TestCreateTopicIntegration(t *testing.T) {
	t.Run("inserts topic into DB and returns topic object", func(t *testing.T) {
		th := testSetup()
		title := "Topic I've added"

		store := New(th.DB)
		topic, _ := store.CreateTopic(title)
		topicID := topic.ID
		store.CreateMessage(&models.Message{TopicID: topicID, Content: "new message", AuthorInitials: "JK", AuthorTheme: 1})
		recentTopics, _ := store.GetRecentTopics()
		firstTopic := recentTopics[0]

		if firstTopic.Title != title {
			t.Error("new topic not present in recent topics list")
		}
	})
}
