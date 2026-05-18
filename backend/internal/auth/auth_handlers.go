package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (o *Orchestration) Login(c *gin.Context) {
	state := NewState()
	SetStateCookie(c, state)
	c.Redirect(http.StatusFound, o.provider.AuthURL(state))
}

func (o *Orchestration) Callback(c *gin.Context) {
	if c.Query("state") != StateCookie(c) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := o.provider.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if !token.IsHuman() {
		c.JSON(http.StatusForbidden, gin.H{"error": "biometric verification required"})
		return
	}

	user, err := o.users.Upsert(token.Sub())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	sessionToken, err := o.sessions.Create(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	SetSessionCookie(c, sessionToken)
	c.Redirect(http.StatusFound, "/")
}

func (o *Orchestration) RequireAuth(c *gin.Context) {
	session, err := o.sessions.Get(SessionCookie(c))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	SetCurrentUser(c, session.UserID)
	c.Next()
}

func (o *Orchestration) Logout(c *gin.Context) {
	o.sessions.Delete(SessionCookie(c))
	ClearSessionCookie(c)
	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}