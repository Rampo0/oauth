package rest

import (
	"encoding/json"
	"os"
	"time"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/rampo0/go-utils/rest_error"
	"github.com/rampo0/multi-lang-microservice/oauth/src/domain/users"
)

const (
	usersService = "USERS_SERVICE"
)

var (
	userRestClient = rest.RequestBuilder{
		BaseURL: os.Getenv(usersService),
		Timeout: 500 * time.Millisecond,
	}
)

type UserRepository interface {
	LoginUser(string, string) (*users.User, *rest_error.RestErr)
}

type userRepository struct {
}

func NewRestUserRepository() *userRepository {
	return &userRepository{}
}

func (r *userRepository) LoginUser(email string, password string) (*users.User, *rest_error.RestErr) {

	request := users.LoginRequest{
		Email:    email,
		Password: password,
	}

	response := userRestClient.Post("/users/login", request)

	if response == nil || response.Response == nil {
		return nil, rest_error.NewInternalServerError(os.Getenv(usersService))
	}

	if response.StatusCode > 299 {
		var restErr rest_error.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, rest_error.NewInternalServerError("invalid error interface when trying to login user")
		}

		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, rest_error.NewInternalServerError("error when trying unmarshall users response")
	}

	return &user, nil

}
