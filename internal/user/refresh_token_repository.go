package user

import "gorm.io/gorm"

type RefreshTokenRepository interface {
	Create(rt *RefreshToken) error
	FindByToken(token string) (*RefreshToken, error)
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(rt *RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *refreshTokenRepository) FindByToken(token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	return &refreshToken, err
}
