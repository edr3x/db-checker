package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	pg_err := false
	redis_err := false
	pg_message := "Postgres connection successful"
	redis_message := "Redis connection successful"

	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		pg_err = true
		log.Println("Database connection failed:", err)
		pg_message = fmt.Sprintf("PG connection failed as %s", err)
	}
	defer conn.Close(ctx)

	url := os.Getenv("REDIS_URL")
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Println("Error connecting to redis")
		redis_err = true
		redis_message = fmt.Sprintf("Redis connection failed as %s", err)
	}
	redis.NewClient(opts)

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		generated_message := generatePage(pg_err, redis_err, pg_message, redis_message)
		if pg_err || redis_err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(generated_message))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(generated_message))
	})

	log.Println("Listening on port 8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func generatePage(pg_err bool, redis_err bool, pg_message string, redis_message string) string {
	if pg_err || redis_err {
		return fmt.Sprintf("<html><body><h1>%s</h1><h1>%s</h1></body></html>", pg_message, redis_message)
	}
	return fmt.Sprintf("<html><body><h1>%s</h1><h1>%s</h1></body></html>", pg_message, redis_message)
}
