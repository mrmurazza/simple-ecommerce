package impl

import (
	"ecommerce/domain/user"
	"ecommerce/dto/request"
	"ecommerce/pkg/auth"

	"errors"
)

type service struct {
	repo    user.Repository
	authSvc auth.Service
}

func NewService(repo user.Repository, authSvc auth.Service) user.Service {
	return &service{
		repo:    repo,
		authSvc: authSvc,
	}
}

func (s *service) ConvertToRole(roleStr string) (user.Role, error) {
	roleMap := user.GetAllRoles()

	role := user.Role(roleStr)
	if _, ok := roleMap[role]; !ok {
		return "", errors.New("role does not exist")
	}

	return role, nil
}

func (s *service) CreateUserIfNotAny(req request.CreateUserRequest) (*user.User, error) {
	role, err := s.ConvertToRole(req.Role)
	if err != nil {
		return nil, err
	}

	password := s.authSvc.GeneratePassword(4)
	u := &user.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: s.authSvc.EncryptPassword(password),
		Role:     role,
	}

	existingUser, err := s.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this phonenumber already exist")
	}
	if err != nil {
		return nil, err
	}

	u, err = s.repo.Persist(u)
	if err != nil {
		return nil, err
	}

	// overwrite encrypted password for response purposes
	u.Password = password

	return u, nil
}

func (s *service) Login(email, password string) (*user.User, string, error) {

	u, err := s.repo.GetUserByUserPass(email, s.authSvc.EncryptPassword(password))
	if u == nil {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}

	claims := map[string]interface{}{
		"email":     u.Email,
		"name":      u.Name,
		"role":      u.Role,
		"timestamp": u.CreatedAt.UTC().Unix(),
	}
	token, err := s.authSvc.TokenizeData(claims)
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}
