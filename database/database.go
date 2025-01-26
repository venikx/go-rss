package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/venikx/go-rss/typings"
)

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
	schema   = os.Getenv("DB_SCHEMA")

	db *pgxpool.Pool
)

func New() *pgxpool.Pool {
	if db != nil {
		return db
	}

	ctx := context.Background()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse connection string of database: %v\n", err)
		os.Exit(1)
	}

	db, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create pool of database: %v\n", err)
		os.Exit(1)
	}

	return db
}

func Close() {
	db.Close() // Close the pool to release connections when the app shuts down
}

func ReadUsers(ctx context.Context) ([]typings.User, error) {
	query := "SELECT id, name, api_key FROM users"
	users := make([]typings.User, 0)

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := typings.User{}
		_ = rows.Scan(&user.Id, &user.Name, &user.ApiKey)

		users = append(users, user)
	}

	return users, err
}

func CreateUser(ctx context.Context, name string) (typings.User, error) {
	query := "INSERT INTO users(name) VALUES (@userName) RETURNING id, name, api_key"
	args := pgx.NamedArgs{
		"userName": name,
	}
	var user = typings.User{}

	err := db.QueryRow(ctx, query, args).Scan(&user.Id, &user.Name, &user.ApiKey)
	return user, err
}

func ReadFeeds(ctx context.Context) ([]typings.Feed, error) {
	query := "SELECT id, name, url, user_id FROM feeds"
	feeds := make([]typings.Feed, 0)

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		feed := typings.Feed{}
		_ = rows.Scan(&feed.Id, &feed.Name, &feed.Url, &feed.UserId)

		feeds = append(feeds, feed)
	}

	return feeds, err
}

func CreateFeed(ctx context.Context, name string, url string, userId int) (typings.Feed, error) {
	query := "INSERT INTO feeds(name, url, user_id) VALUES ($1, $2, $3) RETURNING id, name, url, user_id"
	var feed = typings.Feed{}

	err := db.QueryRow(ctx, query, name, url, userId).Scan(&feed.Id, &feed.Name, &feed.Url, &feed.UserId)
	return feed, err
}

func FollowFeed(ctx context.Context, feedId int, userId int) (typings.FeedFollow, error) {
	query := "INSERT INTO feed_follows(feed_id, user_id) VALUES ($1, $2) RETURNING id"
	var feed = typings.FeedFollow{}

	err := db.QueryRow(ctx, query, feedId, userId).Scan(&feed.Id)
	return feed, err
}

//func (s *service) GetFeedFollows(ctx context.Context, userId int) (typings.FeedFollow, error) {
//	query := "SELECT * from feed_follows WHERE user_id=$1 RETURNING id, feed_id"
//	var feed = typings.FeedFollow{}
//
//	err := s.db.QueryRow(ctx, query, userId).Scan(&feed.Id, &feed.FeedId)
//	return feed, err
//}

func GetNextFeedsToFetch(ctx context.Context, limit int) ([]typings.Feed, error) {
	query := "SELECT id, name, url, user_id FROM feeds ORDER BY last_fetched_at ASC NULLS FIRST LIMIT $1"
	feeds := make([]typings.Feed, 0)

	rows, err := db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		feed := typings.Feed{}
		if err := rows.Scan(&feed.Id, &feed.Name, &feed.Url, &feed.UserId); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %w", err)
	}

	return feeds, nil
}

func MarkFeedFetched(ctx context.Context, feedId int) (typings.Feed, error) {
	query := "UPDATE feeds SET last_fetched_at = NOW() WHERE id=$1 RETURNING id, name, url, user_id"
	var feed = typings.Feed{}

	err := db.QueryRow(ctx, query, feedId).Scan(&feed.Id, &feed.Name, &feed.Url, &feed.UserId)
	return feed, err
}
