package web

import (
	"zebra-rss/storage"

	"github.com/gin-gonic/gin"
)

func StartServer(storage *storage.Storage) {
	handler := handler{storage}

	router := gin.Default()
	router.Use(SessionMiddleware)
	{
		router.POST("/signin", handler.signIn)
		router.POST("/logout", handler.signOut)

		safe := router.Group("")
		safe.Use(AuthenticatedMiddleware)
		{
			safe.GET("/me", handler.user)
		}
	}

	router.Run("0.0.0.0:8080")
}
