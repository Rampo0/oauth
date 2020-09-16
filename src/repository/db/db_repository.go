package db

import (
	"github.com/gocql/gocql"
	"github.com/rampo0/go-utils/rest_error"
	"github.com/rampo0/multi-lang-microservice/oauth/src/clients/cassandra"
	"github.com/rampo0/multi-lang-microservice/oauth/src/domain/access_token"
)

const (
	queryGetAccessToken    = "SELECT access_token , user_id , client_id, expires FROM access_token WHERE access_token=?;"
	queryCreateAccessToken = "INSERT INTO access_token(access_token, user_id, client_id, expires) VALUES (?,?,?,?);"
	queryUpdateExpires     = "UPDATE access_token SET expires=? WHERE access_token=?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *rest_error.RestErr)
	Create(access_token.AccessToken) *rest_error.RestErr
	UpdateExpires(access_token.AccessToken) *rest_error.RestErr
}

type dbRepository struct {
}

func NewRepository() *dbRepository {
	return &dbRepository{}
}

func (repo *dbRepository) UpdateExpires(at access_token.AccessToken) *rest_error.RestErr {

	if err := cassandra.GetSession().Query(queryUpdateExpires,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return rest_error.NewInternalServerError("error when trying to Update expires in database")
	}

	return nil
}

func (repo *dbRepository) Create(at access_token.AccessToken) *rest_error.RestErr {

	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessToken,
		at.UserID,
		at.ClientID,
		at.Expires,
	).Exec(); err != nil {
		return rest_error.NewInternalServerError("error when trying to save access token in database")
	}

	return nil
}

func (repo *dbRepository) GetById(id string) (*access_token.AccessToken, *rest_error.RestErr) {

	var accessToken access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&accessToken.AccessToken,
		&accessToken.UserID,
		&accessToken.ClientID,
		&accessToken.Expires,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_error.NewNotFoundError("No access token found with given id")
		}

		return nil, rest_error.NewInternalServerError(err.Error())
	}

	return &accessToken, nil
}
