package web

import (
	"zebra-rss/storage"

	"github.com/gin-gonic/gin"
)

func StartServer(storage *storage.Storage) {
	handler := handler{storage}

	router := gin.Default()
	router.Use(SessionMiddleware)
	router.POST("/signin", handler.signIn)
	router.POST("/logout", handler.signOut)
	router.Run("localhost:8080")
}
