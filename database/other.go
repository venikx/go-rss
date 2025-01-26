package database

import (
	"context"
	"fmt"
	"log"
)

func Health(ctx context.Context) map[string]string {
	stats := make(map[string]string)

	err := db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy!"

	// TODO(Kevin): Connect to pool to view stats
	//dbStats := s.pgxpool.Stats()
	//stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	//stats["in_use"] = strconv.Itoa(dbStats.InUse)
	//stats["idle"] = strconv.Itoa(dbStats.Idle)
	//stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	//stats["wait_duration"] = dbStats.WaitDuration.String()
	//stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	//stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	//if dbStats.OpenConnections > 40 {
	//	stats["message"] = "The database is experiencing heavy load."
	//}

	//if dbStats.WaitCount > 1000 {
	//	stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	//}

	//if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
	//	stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	//}

	//if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
	//	stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	//}

	return stats
}

func HelloWorld(ctx context.Context) (string, error) {
	var greeting string

	err := db.QueryRow(ctx, "SELECT 'Hello, world!'").Scan(&greeting)
	if err != nil {
		return "", err
	}

	return greeting, nil
}
