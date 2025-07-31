package web

import (
	"zebra-rss/storage"

	"github.com/gin-gonic/gin"
)

type handler struct {
	storage *storage.Storage
}

func (h *handler) signIn(c *gin.Context) {
	data := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.ShouldBindBodyWithJSON(&data)

	if err != nil {
		c.JSON(400, gin.H{"error": "invalid body format"})

		return
	}

	if h.storage.Users.VerifyCredentials(data.Email, data.Password) != nil {
		c.JSON(401, gin.H{"error": "wrong credentials"})

		return
	}

	user, err := h.storage.Users.GetUserByEmail(data.Email)

	if err != nil {
		c.JSON(401, gin.H{"error": "wrong credentials"})

		return
	}

	// Close previous session.
	if session := GinGetSession(c); session != nil {
		session.Close()
	}

	// Start new session.
	session, err := GinStartNewSession(c)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to start a new session"})

		return
	}

	session.SetInt64("user_id", user.Id)

	// Return the logged in user.
	c.JSON(200, user)
}

func (h *handler) signOut(c *gin.Context) {
	// Close previous session.
	if session := GinGetSession(c); session != nil {
		session.Close()
	}

	// Start new session.
	if _, err := GinStartNewSession(c); err != nil {
		c.JSON(500, gin.H{"error": "failed to start a new session"})

		return
	}

	c.String(200, "")
}
