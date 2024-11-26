package repositories

import (
	"online-music-library/internal/models"

	"gorm.io/gorm"
)

type SongRepository struct {
	DB *gorm.DB
}

func (repo *SongRepository) Add(song *models.Song) error {
	return repo.DB.Create(song).Error
}
func (repo *SongRepository) Delete(id uint) error {
	return repo.DB.Delete(&models.Song{}, id).Error
}

func (repo *SongRepository) Update(song *models.Song) error {
	return repo.DB.Save(song).Error
}

func (repo *SongRepository) GetAll(filters map[string]string, offset, limit int) ([]models.Song, error) {
	var songs []models.Song
	query := repo.DB

	if filters["group"] != "" {
		query = query.Where("group = ?", filters["group"])
	}

	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}
