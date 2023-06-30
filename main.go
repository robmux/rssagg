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

	dbQueries := database.New(db)
	apiconfig := apiConfig{DB: dbQueries}

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

	v1Router.Post("/users", apiconfig.handlerCreateUser)
	v1Router.Get("/users", apiconfig.middlewareAuth(apiconfig.handlerGetUser))

	// Feeds
	v1Router.Post("/feeds", apiconfig.middlewareAuth(apiconfig.handlerCreateFeed))
	v1Router.Get("/feeds", apiconfig.handlerGetFeeds)

	// Feed Follows
	v1Router.Post("/feed_follows", apiconfig.middlewareAuth(apiconfig.handlerCreateFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiconfig.middlewareAuth(apiconfig.handlerDeleteFeedFollow))
	v1Router.Get("/feed_follows", apiconfig.middlewareAuth(apiconfig.handlerGetFeedFollows))

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
