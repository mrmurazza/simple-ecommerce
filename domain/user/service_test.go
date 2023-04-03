package user_test

import (
	"ecommerce/domain/user"
	"ecommerce/domain/user/impl"
	"ecommerce/domain/user/mocks"
	"ecommerce/dto/request"
	authMocks "ecommerce/pkg/auth/mocks"
	"errors"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
)

type userServiceTest struct {
	MockRepo    *mocks.Repository
	MockService *mocks.Service
	MockAuth    *authMocks.Service
	Service     user.Service
}

var svcTest userServiceTest

func init() {
	mockRepo := new(mocks.Repository)
	mockService := new(mocks.Service)
	mockAuth := new(authMocks.Service)

	svcTest = userServiceTest{
		MockRepo:    mockRepo,
		MockService: mockService,
		MockAuth:    mockAuth,
		Service: impl.NewService(
			mockRepo,
			mockAuth,
		),
	}
}

func TestCreateUserIfNotAny(test *testing.T) {
	req := request.CreateUserRequest{
		Email: "123456",
		Name:  "User-123",
		Role:  "ADMIN",
	}

	u := &user.User{
		Name:     "User-123",
		Email:    "123456",
		Password: "ASDF",
		Role:     user.RoleAdmin,
	}

	userWithId := &user.User{
		ID:       1,
		Name:     "User-123",
		Email:    "123456",
		Password: "QWER",
		Role:     user.RoleAdmin,
	}

	test.Run("error role not valid", func(t *testing.T) {
		req := request.CreateUserRequest{}
		expectedErr := errors.New("role does not exist")

		res, err := svcTest.Service.CreateUserIfNotAny(req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, expectedErr, err)
	})

	test.Run("error Database", func(t *testing.T) {
		expectedErr := errors.New("err Database")

		svcTest.MockAuth.On("GeneratePassword", 4).Return("QWER").Once()
		svcTest.MockAuth.On("EncryptPassword", "QWER").Return("ASDF").Once()
		svcTest.MockRepo.On("GetUserByPhonenumber", "123456").Return(nil, expectedErr).Once()

		resAct, resErr := svcTest.Service.CreateUserIfNotAny(req)

		assert.True(t, svcTest.MockRepo.AssertNotCalled(t, "Persist", u))
		assert.True(t, svcTest.MockAuth.AssertExpectations(t), "mock method from mock auth not called as expected")
		assert.True(t, svcTest.MockRepo.AssertExpectations(t), "mock method from mock repo not called as expected")
		assert.Nil(t, resAct)
		assert.Equal(t, expectedErr, resErr)

	})

	test.Run("error users already exist", func(t *testing.T) {
		svcTest.MockAuth.On("GeneratePassword", 4).Return("QWER").Once()
		svcTest.MockAuth.On("EncryptPassword", "QWER").Return("ASDF").Once()
		svcTest.MockRepo.On("GetUserByPhonenumber", "123456").Return(userWithId, nil).Once()

		expectedErr := errors.New("user with this phonenumber already exist")
		resAct, resErr := svcTest.Service.CreateUserIfNotAny(req)

		assert.True(t, svcTest.MockRepo.AssertNotCalled(t, "Persist", u))
		assert.True(t, svcTest.MockAuth.AssertExpectations(t), "mock method from mock auth not called as expected")
		assert.True(t, svcTest.MockRepo.AssertExpectations(t), "mock method from mock repo not called as expected")
		assert.Nil(t, resAct)
		assert.Equal(t, expectedErr, resErr)

	})

	test.Run("positive result", func(t *testing.T) {
		svcTest.MockAuth.On("GeneratePassword", 4).Return("QWER").Once()
		svcTest.MockAuth.On("EncryptPassword", "QWER").Return("ASDF").Once()
		svcTest.MockRepo.On("GetUserByPhonenumber", "123456").Return(nil, nil).Once()
		svcTest.MockRepo.On("Persist", u).Return(userWithId, nil).Once()
		res, err := svcTest.Service.CreateUserIfNotAny(req)

		assert.True(t, svcTest.MockAuth.AssertExpectations(t), "mock method from mock auth not called as expected")
		assert.True(t, svcTest.MockRepo.AssertExpectations(t), "mock method from mock repo not called as expected")
		assert.Nil(t, err)
		assert.Equal(t, userWithId, res)

	})
}

func TestLogin(test *testing.T) {
	phonenumber := "123456"
	password := "QWER"
	encryptedPass := "ASDF"
	ts := int64(123123132)
	createdAt := time.Unix(ts, 0)

	u := &user.User{
		Email:     phonenumber,
		Name:      "USER-123",
		Role:      user.RoleAdmin,
		CreatedAt: &createdAt,
	}

	claims := map[string]interface{}{
		"phonenumber": phonenumber,
		"name":        "USER-123",
		"role":        user.RoleAdmin,
		"timestamp":   ts,
	}

	token := "TOKEN-123"

	test.Run("login user not found", func(t *testing.T) {
		svcTest.MockAuth.On("EncryptPassword", password).Return(encryptedPass).Once()
		svcTest.MockRepo.On("GetUserByUserPass", phonenumber, encryptedPass).Return(nil, nil).Once()

		resAct, tokenAct, errAct := svcTest.Service.Login(phonenumber, password)

		assert.Nil(t, resAct)
		assert.Nil(t, errAct)
		assert.Empty(t, tokenAct)
	})

	test.Run("error on tokenize data", func(t *testing.T) {
		expectedErr := errors.New("err tokenize")

		svcTest.MockAuth.On("EncryptPassword", password).Return(encryptedPass).Once()
		svcTest.MockRepo.On("GetUserByUserPass", phonenumber, encryptedPass).Return(u, nil).Once()
		svcTest.MockAuth.On("TokenizeData", claims).Return("", expectedErr).Once()

		resAct, tokenAct, errAct := svcTest.Service.Login(phonenumber, password)

		assert.True(t, svcTest.MockAuth.AssertExpectations(t), "mock method from mock auth not called as expected")
		assert.True(t, svcTest.MockRepo.AssertExpectations(t), "mock method from mock repo not called as expected")
		assert.Nil(t, resAct)
		assert.Empty(t, tokenAct)
		assert.Equal(t, expectedErr, errAct)

	})

	test.Run("positive result", func(t *testing.T) {
		svcTest.MockAuth.On("EncryptPassword", password).Return(encryptedPass).Once()
		svcTest.MockRepo.On("GetUserByUserPass", phonenumber, encryptedPass).Return(u, nil).Once()
		svcTest.MockAuth.On("TokenizeData", claims).Return(token, nil).Once()
		resAct, tokenAct, errAct := svcTest.Service.Login(phonenumber, password)

		assert.True(t, svcTest.MockAuth.AssertExpectations(t), "mock method from mock auth not called as expected")
		assert.True(t, svcTest.MockRepo.AssertExpectations(t), "mock method from mock repo not called as expected")
		assert.Equal(t, u, resAct)
		assert.Equal(t, token, tokenAct)
		assert.Nil(t, errAct)
	})
}
