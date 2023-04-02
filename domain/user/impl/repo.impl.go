package impl

import (
	"ecommerce/domain/user"

	"github.com/jinzhu/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) user.Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) Persist(u *user.User) (*user.User, error) {
	res := r.db.Save(u)

	err := res.Error
	if err != nil {
		return nil, err
	}

	res.Last(u)

	return u, nil

}

func (r *repo) GetUserByUserPass(phonenumber, password string) (*user.User, error) {
	u := user.User{}
	err := r.db.Model(user.User{}).Where("phonenumber = ? AND password = ?", phonenumber, password).First(&u).Error
	if err != nil {
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return &u, nil
}

func (r *repo) GetUserByPhonenumber(phonenumber string) (*user.User, error) {
	u := user.User{}
	err := r.db.Model(user.User{}).Where("phonenumber = ?", phonenumber).First(&u).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}
