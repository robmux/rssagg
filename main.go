package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/robmux/rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening db connection %v", err)
	}

	// Link sqlc queries with the db connection
	dbQueries := database.New(db)
	apiConfig := apiConfig{DB: dbQueries}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatalln("PORT is not found in the environment")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	v1Router.Post("/users", apiConfig.handlerCreateUser)
	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handlerGetUser))

	// Feeds
	v1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	v1Router.Get("/feeds", apiConfig.handlerGetFeeds)

	// Feed Follows
	v1Router.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerCreateFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.handlerDeleteFeedFollow))
	v1Router.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerGetFeedFollows))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}

	fmt.Println("Server starting on PORT: ", portStr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln("Error starting server, ", err)
	}
}
