package request

import "errors"

type (
	CreateUserRequest struct {
		Phonenumber string `json:"phonenumber"`
		Name        string `json:"name"`
		Role        string `json:"role"`
	}

	LoginRequest struct {
		Phonenumber string `json:"phonenumber"`
		Password    string `json:"password"`
	}
)

func (cr CreateUserRequest) Validate() error {
	if cr.Phonenumber == "" {
		return errors.New("phonenumber is required")
	}

	if cr.Role == "" {
		return errors.New("role is required")
	}

	return nil
}
