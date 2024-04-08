package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	ReadUsers() ([]User, error)
	CreateUser(name string) (User, error)
	ReadFeeds() ([]Feed, error)
	CreateFeed(name string, url string, userId int) (Feed, error)
	Health() (map[string]string, error)
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}
	return s
}

type User struct {
	Id     int
	Name   string
	ApiKey string
}

func (s *service) ReadUsers() ([]User, error) {
	query := "SELECT id, name, api_key FROM users"
	users := make([]User, 0)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := User{}
		_ = rows.Scan(&user.Id, &user.Name, &user.ApiKey)

		users = append(users, user)
	}

	return users, err
}

func (s *service) CreateUser(name string) (User, error) {
	query := "INSERT INTO users(name) VALUES (@userName) RETURNING id, name, api_key"
	args := pgx.NamedArgs{
		"userName": name,
	}
	var user = User{}

	err := s.db.QueryRow(query, args).Scan(&user.Id, &user.Name, &user.ApiKey)
	return user, err
}

type Feed struct {
	Id     int
	Name   string
	Url    string
	UserId int
}

func (s *service) ReadFeeds() ([]Feed, error) {
	query := "SELECT id, name, url, user_id FROM feeds"
	feeds := make([]Feed, 0)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		feed := Feed{}
		_ = rows.Scan(&feed.Id, &feed.Name, &feed.Url, &feed.UserId)

		feeds = append(feeds, feed)
	}

	return feeds, err
}

func (s *service) CreateFeed(name string, url string, userId int) (Feed, error) {
	query := "INSERT INTO feeds(name, url, user_id) VALUES ($1, $2, $3) RETURNING id, name, url, user_id"
	var feed = Feed{}

	err := s.db.QueryRow(query, name, url, userId).Scan(&feed.Id, &feed.Name, &feed.Url, &feed.UserId)
	return feed, err
}

func (s *service) Health() (map[string]string, error) {
	err := s.db.Ping()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"message": "It's healthy",
	}, err
}
