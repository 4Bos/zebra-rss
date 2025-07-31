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
