package impl

import (
	"ecommerce/domain/user"

	"gorm.io/gorm"
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

func (r *repo) GetUserByUserPass(email, password string) (*user.User, error) {
	u := user.User{}
	err := r.db.Model(user.User{}).Where("email = ? AND password = ?", email, password).First(&u).Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *repo) GetUserByEmail(email string) (*user.User, error) {
	u := user.User{}
	err := r.db.Model(user.User{}).Where("email = ?", email).First(&u).Error

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *repo) GetByIds(ids []int) ([]user.User, error) {
	var users []user.User

	q := r.db.Model(user.User{}).
		Where("id in (?)", ids).
		Find(&users)

	err := q.Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
