package request

import (
	"errors"
	"fmt"
)

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

	CreateOrderRequest struct {
		UserId   int                         `json:"user_id"`
		Products []*CreateOrderRequestDetail `json:"products"`
	}

	CreateOrderRequestDetail struct {
		ProductId int `json:"product_id"`
		Qty       int `json:"qty"`
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

func (co CreateOrderRequest) Validate() error {
	if co.UserId == 0 {
		return errors.New("UserId is required")
	}

	if len(co.Products) == 0 {
		return errors.New("Must select at least one product")
	}

	for idx, p := range co.Products {
		err := p.Validate()
		if err != nil {
			return fmt.Errorf("found error on product idx %d: %s", idx, err.Error())
		}
	}

	return nil
}

func (co CreateOrderRequestDetail) Validate() error {
	if co.ProductId == 0 {
		return errors.New("ProductId is required")
	}

	if co.Qty == 0 {
		return errors.New("Qty must be greater than 0")
	}

	return nil
}
