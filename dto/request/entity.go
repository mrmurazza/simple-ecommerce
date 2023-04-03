package request

import "errors"

type (
	CreateUserRequest struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Role  string `json:"role"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

func (cr CreateUserRequest) Validate() error {
	if cr.Name == "" {
		return errors.New("Name is required")
	}

	if cr.Email == "" {
		return errors.New("Email is required")
	}

	if cr.Role == "" {
		return errors.New("role is required")
	}

	return nil
}
