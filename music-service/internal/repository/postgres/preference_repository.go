package postgres

import "gorm.io/gorm"

type PreferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) *PreferenceRepository {
	return &PreferenceRepository{db}
}

func (r *PreferenceRepository) UpsertUserGenreScore(userID, genre string, score int) error {
	return r.db.
		Where("user_id = ? AND genre = ?", userID, genre).
		Assign(UserGenreScore{Score: score}).
		FirstOrCreate(&UserGenreScore{
			UserID: userID,
			Genre:  genre,
		}).Error
}
