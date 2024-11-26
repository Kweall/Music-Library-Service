package main

import (
	"log"
	"net/http"
	"online-music-library/internal/handlers"
	"online-music-library/internal/repositories"
	"online-music-library/internal/services"
	"online-music-library/internal/utils"
	"online-music-library/migrations"

	"github.com/gorilla/mux"
)

func main() {
	db := utils.InitDB()
	migrations.Migrate(db)

	repo := &repositories.SongRepository{DB: db}
	service := &services.SongService{Repo: repo}
	handler := &handlers.SongHandler{Service: service}

	r := mux.NewRouter()

	r.HandleFunc("/songs", handler.GetSongs).Methods("GET")
	r.HandleFunc("/songs", handler.AddSong).Methods("POST")
	r.HandleFunc("/songs/{id:[0-9]+}", handler.DeleteSong).Methods("DELETE")
	r.HandleFunc("/songs/{id:[0-9]+}", handler.UpdateSong).Methods("PUT")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
