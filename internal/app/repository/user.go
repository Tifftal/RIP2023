package repository

import "MSRM/internal/app/ds"

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
