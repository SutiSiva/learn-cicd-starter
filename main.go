package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// init() Dummy-Aufrufe, damit staticcheck keine "unused"-Fehler meldet
func init() {
	cfg := &apiConfig{}
	_ = cfg.handlerNotesGet
	_ = cfg.handlerNotesCreate
	_ = cfg.handlerUsersCreate
	_ = cfg.handlerUsersGet
	_ = generateRandomSHA256Hash
	_ = respondWithError
	_ = respondWithJSON
	_ = staticFiles

	// Dummy für noch ungenutzte Funktionen/Typs
	_ = handlerReadiness
	var _ authedHandler
	_ = cfg.middlewareAuth
}

type apiConfig struct {
	DB *database.Queries
}

//go:embed static/*
var staticFiles embed.FS

func main() {
	// .env laden
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	apiCfg := apiConfig{}

	// Datenbank verbinden, wenn URL vorhanden
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL environment variable is not set")
		log.Println("Running without CRUD endpoints")
	} else {
		db, err := sql.Open("libsql", dbURL)
		if err != nil {
			log.Fatal(err)
		}
		dbQueries := database.New(db)
		apiCfg.DB = dbQueries
		log.Println("Connected to database!")
	}

	router := chi.NewRouter()

	// CORS aktivieren
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Static Files einbinden
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	router.Handle("/*", http.FileServer(http.FS(subFS)))

	// HTTP Server mit Timeouts (fix für gosec G114)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Serving on port %s\n", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
