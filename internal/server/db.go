package server

import (
	"context"
	"errors"
	"sync"
	"time"

	"hirelin-auth/internal/logger"
	database "hirelin-auth/internal/server/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db   *database.Queries
	pool *pgxpool.Pool
	mu   sync.RWMutex
)

func SetDB(d *database.Queries) {
	mu.Lock()
	defer mu.Unlock()
	db = d
}

func GetDB() *database.Queries {
	mu.RLock()
	defer mu.RUnlock()
	return db
}

func ConnectDB() (*database.Queries, error) {
	config, err := pgxpool.ParseConfig(logger.GetEnv("DATABASE_URL"))

	if err != nil {
		return nil, errors.New("invalid database configuration")
	}

	config.MaxConns = 5
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	poolConn, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return nil, errors.New("failed to connect to database")
	}

	if err := poolConn.Ping(context.Background()); err != nil {
		return nil, errors.New("couldn't reach database")
	}

	logger.Logger.Println(logger.Colors["blue"] + "Connected to database" + logger.Colors["reset"])

	pool = poolConn
	queries := database.New(pool)
	SetDB(queries)

	go monitorDatabaseConnection()

	return queries, nil
}

func ReconnectDB() (*database.Queries, error) {
	mu.Lock()
	defer mu.Unlock()

	logger.Logger.Println(logger.Colors["yellow"] + "Attempting to reconnect to database..." + logger.Colors["reset"])

	if pool != nil {
		pool.Close()
	}

	return ConnectDB()
}

func IsDatabaseConnected() bool {
	if pool == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := pool.Ping(ctx)
	return err == nil
}

func monitorDatabaseConnection() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !IsDatabaseConnected() {
			logger.Logger.Println(logger.Colors["red"] + "Database connection lost!" + logger.Colors["reset"])
			_, err := ReconnectDB()
			if err != nil {
				logger.Logger.Printf("%sFailed to reconnect: %v%s\n",
					logger.Colors["red"], err, logger.Colors["reset"])
			} else {
				logger.Logger.Println(logger.Colors["green"] + "Successfully reconnected to database" + logger.Colors["reset"])
			}
		}
	}
}
