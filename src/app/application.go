package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rampo0/multi-lang-microservice/oauth/src/http"
	"github.com/rampo0/multi-lang-microservice/oauth/src/repository/db"
	"github.com/rampo0/multi-lang-microservice/oauth/src/repository/rest"
	"github.com/rampo0/multi-lang-microservice/oauth/src/services/access_token"
)

var (
	router = gin.Default()
)

func StartApplication() {

	// if cassandra.GetSession() == nil {
	// 	panic("Unable to connect to cassandra db")
	// }

	atHandler := http.NewAccessTokenHandler(
		access_token.NewService(db.NewRepository(), rest.NewRestUserRepository()))

	router.GET("/oauth/access-token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access-token", atHandler.Create)

	router.Run(":8081")
}
