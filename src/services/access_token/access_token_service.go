package access_token

import (
	"strings"

	"github.com/rampo0/go-utils/rest_error"
	"github.com/rampo0/multi-lang-microservice/oauth/src/domain/access_token"

	"github.com/rampo0/multi-lang-microservice/oauth/src/repository/db"

	"github.com/rampo0/multi-lang-microservice/oauth/src/repository/rest"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, *rest_error.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *rest_error.RestErr)
	UpdateExpires(access_token.AccessToken) *rest_error.RestErr
}

type service struct {
	repository    db.DbRepository
	restUsersRepo rest.UserRepository
}

func NewService(repo db.DbRepository, restRepo rest.UserRepository) *service {
	return &service{
		repository:    repo,
		restUsersRepo: restRepo,
	}
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, *rest_error.RestErr) {

	if err := request.Validate(); err != nil {
		return nil, err
	}

	user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	at := access_token.GetNewAccessToken(user.ID)
	at.Generate()

	if err := s.repository.Create(at); err != nil {
		return nil, err
	}

	return &at, nil
}

func (s *service) UpdateExpires(at access_token.AccessToken) *rest_error.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.repository.UpdateExpires(at)
}

func (s *service) GetById(access_token_id string) (*access_token.AccessToken, *rest_error.RestErr) {

	accessTokenId := strings.TrimSpace(access_token_id)

	if len(accessTokenId) == 0 {
		return nil, rest_error.NewBadRequestError("invalid token id")
	}

	accessToken, err := s.repository.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}
