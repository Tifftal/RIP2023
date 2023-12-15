package repository

import (
	"MSRM/internal/app/ds"
	"context"
	"strconv"
	"time"
)

const blacklistPrefix = "jwt_blacklist."

func (r *Repository) DeleteUserByID(id int) error {
	if err := r.db.Exec("UPDATE users SET user_status='Deleted' WHERE Id_user= ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) EditUser(user *ds.Users) error {
	err := r.db.Where("Id_user = ?", user.Id_user).Updates(&user).Error

	return err
}

func (repository *Repository) GetUserByID(id int) (*ds.Users, error) {
	user := &ds.Users{}

	err := repository.db.First(user, "Id_user = ?", id).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *Repository) GetUserByRole(role string) ([]ds.Users, error) {
	user := []ds.Users{}
	err := repository.db.Where("Role=?", role).Find(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *Repository) GetUserByEmail(email string) (*ds.Users, error) {
	user := &ds.Users{}

	err := repository.db.Where("email_address = ?", email).Find(&user).Error
	if err != nil {
		return &ds.Users{}, err
	}

	return user, nil
}

func (r *Repository) CreateUser(user *ds.Users) error {
	err := r.db.Table("users").Create(&user).Error
	return err
}

func (r *Repository) SaveJWTToken(id uint, token string) error {
	expiration := 1 * time.Hour

	idStr := strconv.FormatUint(uint64(id), 10)

	err := r.rd.Set(idStr, token, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddTokenToBlacklist(id uint, token string, expiration time.Duration) error {
	// Генерируем уникальный ключ для токена в блеклисте
	blacklistKey := blacklistPrefix + token
	idStr := strconv.FormatUint(uint64(id), 10)

	return r.rd.Set(idStr, blacklistKey, expiration).Err()
}

func (r *Repository) IsTokenBlacklisted(ctx context.Context, token string) error {
	// Проверяем наличие токена в блеклисте
	blacklistKey := blacklistPrefix + token

	return r.rd.Get(blacklistKey).Err()
}
