package repositories

import (
	"fmt"

	"github.com/sayopaul/sendchamp-go-test/infrastructure"
	"github.com/sayopaul/sendchamp-go-test/models"
)

type UserRepository struct {
	db infrastructure.Database
}

func NewUserRepository(db infrastructure.Database) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (ur UserRepository) GetOne(field string, value string) (user models.User, err error) {
	query := fmt.Sprintf("%v = ?", field)
	return user, ur.db.DB.Where(query, value).First(&user).Error
}
func (ur UserRepository) FetchUser(userCondition models.User) (user models.User, err error) {
	return user, ur.db.DB.Model(&models.User{}).Where(&userCondition).First(&user).Error
}

func (ur UserRepository) CheckSignIn(email string, password string) (user models.User, error error) {
	return user, ur.db.DB.Where("email = ? AND password = ?", email, password).First(&user).Error
}

func (ur UserRepository) CreateUser(u models.User) (user models.FetchUser, err error) {
	err = ur.db.DB.Create(&u).Error
	ur.db.DB.Model(&models.User{}).Where(models.User{ID: u.ID}).First(&user)
	return user, err
}

func (ur UserRepository) UpdatePassword(user models.User, password string) error {
	return ur.db.DB.Model(&user).Update("password", password).Error
}

func (ur UserRepository) Update(id uint, user models.User) (models.User, error) {
	return user, ur.db.DB.Model(models.User{}).Where("id = ?", id).Updates(&user).Error
}
