package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/rampo0/go-utils/crypto_utils"
	"github.com/rampo0/go-utils/rest_error"
)

const (
	expirationTime       = 24
	grantTypePassword    = "password"
	grantTypeCredentials = "client_credentials"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	ClientID    int64  `json:"client_id,omitempty"`
	Expires     int64  `json:"expires"`
}

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// Used for password grant type
	Username string `json:"username"`
	Password string `json:"password"`

	// Used for client_credentials grant type
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() *rest_error.RestErr {
	switch at.GrantType {
	case grantTypePassword:
		break
	case grantTypeCredentials:
		break
	default:
		return rest_error.NewBadRequestError("invalid grant_type parameter")
	}

	return nil
}

func (at *AccessToken) Validate() *rest_error.RestErr {

	at.AccessToken = strings.TrimSpace(at.AccessToken)

	if at.AccessToken == "" {
		return rest_error.NewBadRequestError("invalid token id")
	}

	if at.UserID <= 0 {
		return rest_error.NewBadRequestError("invalid user id")
	}

	if at.ClientID <= 0 {
		return rest_error.NewBadRequestError("invalid client id")
	}

	if at.Expires <= 0 {
		return rest_error.NewBadRequestError("invalid expires")
	}

	return nil
}

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserID:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	now := time.Now().UTC()
	expirationTime := time.Unix(at.Expires, 0)
	return expirationTime.Before(now)
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMD5(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}
