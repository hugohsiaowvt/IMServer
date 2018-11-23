package server

import (

	"rz/config"
	"service/relationservice"
	"service/smsservice"
	"service/userservice"
	"service/utilservice"

	"github.com/gin-gonic/gin"

)

func StartHttpServer() error {
	router := gin.Default()
	setV1Group(router)
	router.Run(":" + config.HTTP_PORT)
	return nil
}

func setV1Group(router gin.IRouter) {

	v := router.Group("/v1")
	{

		users := v.Group("/users")
		{
			users.GET("/test", userservice.Test)


			users.POST("/login/sessions", userservice.Login)
			users.PATCH("update/password", userservice.SetPassword)
			register := users.Group("/register")
			{
				register.POST("/code/:code", userservice.CreateUser)
			}

		}

		sms := v.Group("/sms")
		{
			sms.GET("", smsservice.GetVerificationCode)
			sms.POST("verification/code/:code", smsservice.PreCheckVerificationCode)
		}

		relations := v.Group("/relations")
		{
			relations.GET("/sync/:time/user/:userId", relationservice.SyncFriendList)
			relations.POST("/invite/user/:fromId", relationservice.InviteFriend)
			relations.POST("/accept/user/:fromId", relationservice.AcceptFriend)
		}

		utils := v.Group("/utils")
		{
			utils.GET("/countrys", utilservice.GetCountrys)
		}

	}
}