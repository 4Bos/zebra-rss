package web

import (
	"net/http"
	"net/url"
	"zebra-rss/sessions"

	"github.com/gin-gonic/gin"
)

func SessionMiddleware(c *gin.Context) {
	cookie, err := c.Request.Cookie("session")

	if err != http.ErrNoCookie {
		sessionId, err := url.QueryUnescape(cookie.Value)

		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid session cookie value"})

			return
		}

		session := sessions.GetSession(sessionId)

		if session == nil {
			session, _ = sessions.StartSession()
		}

		c.Set("session", session)
	}

	c.Next()
}

func AuthenticatedMiddleware(c *gin.Context) {
	session := GinGetSession(c)

	if session == nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthenticated"})

		return
	}

	_, exists := session.GetInt64("user_id")

	if !exists {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthenticated"})

		return
	}

	c.Next()
}

func GinGetUserId(c *gin.Context) int64 {
	session := GinGetSession(c)

	if session == nil {
		return 0
	}

	userId, exists := session.GetInt64("user_id")

	if !exists {
		return 0
	}

	return userId
}

func GinGetSession(c *gin.Context) *sessions.Session {
	raw, exists := c.Get("session")

	if !exists {
		return nil
	}

	session, ok := raw.(*sessions.Session)

	if !ok {
		return nil
	}

	return session
}

func GinStartNewSession(c *gin.Context) (*sessions.Session, error) {
	session, err := sessions.StartSession()

	if err != nil {
		return nil, err
	}

	c.Set("session", session)
	c.SetCookie("session", session.Id, 3600, "/", "localhost", false, true)

	return session, nil
}
