package handlers

import (
	"encoding/json"
	"net/http"
	"online-music-library/internal/models"
	"online-music-library/internal/services"
	"strconv"

	"github.com/gorilla/mux"
)

type SongHandler struct {
	Service *services.SongService
}

// AddSong добавляет песню с запросом к внешнему API
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	detailedSong, err := h.Service.FetchSongDetails(song.Group, song.Song)
	if err != nil {
		http.Error(w, "Failed to fetch song details", http.StatusInternalServerError)
		return
	}

	if err := h.Service.AddSong(detailedSong); err != nil {
		http.Error(w, "Failed to add song to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(detailedSong)
}

// DeleteSong удаляет песню по ID
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteSong(uint(id)); err != nil {
		http.Error(w, "Song not found or failed to delete", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateSong обновляет данные песни
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	song.ID = uint(id)

	updatedSong, err := h.Service.UpdateSong(&song)
	if err != nil {
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedSong)
}

func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	filters := map[string]string{
		"group": group,
	}

	// Получаем песни
	songs, err := h.Service.GetSongs(filters, page, size)
	if err != nil {
		http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}
