package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"online-music-library/internal/models"
	"online-music-library/internal/repositories"
	"time"
)

type SongService struct {
	Repo *repositories.SongRepository
}

func (s *SongService) AddSong(song *models.Song) error {
	return s.Repo.Add(song)
}

func (s *SongService) DeleteSong(id uint) error {
	return s.Repo.Delete(id)
}

func (s *SongService) UpdateSong(song *models.Song) (*models.Song, error) {
	err := s.Repo.Update(song)
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (s *SongService) GetSongs(filters map[string]string, page, size int) ([]models.Song, error) {
	offset := (page - 1) * size
	return s.Repo.GetAll(filters, offset, size)
}

func (s *SongService) FetchSongDetails(group, song string) (*models.Song, error) {
	url := "https://pesni.com/info?group=" + group + "&song=" + song

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	var details struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	releaseDate, err := time.Parse("02.01.2006", details.ReleaseDate)
	if err != nil {
		return nil, err
	}

	return &models.Song{
		Group:       group,
		Song:        song,
		ReleaseDate: releaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}, nil
}
