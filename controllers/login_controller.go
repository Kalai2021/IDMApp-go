package controllers

import (
	"html/template"
	"idmapp-go/database"
	"idmapp-go/internal/user"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func loadLoginTemplate() (*template.Template, error) {
	return template.ParseFiles("templates/login.html")
}

func ShowLoginForm(c *gin.Context) {
	redirect := c.Query("redirect")

	// Load template with error handling
	tmpl, err := loadLoginTemplate()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading login template: %v", err)
		return
	}

	c.Header("Content-Type", "text/html")
	c.Status(http.StatusOK)
	if err := tmpl.Execute(c.Writer, gin.H{"redirect": redirect}); err != nil {
		c.String(http.StatusInternalServerError, "Error executing template: %v", err)
	}
}

func HandleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	redirect := c.PostForm("redirect")
	userService := user.NewUserService(database.GetDB())
	user, err := userService.AuthenticateUser(email, password)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		tmpl, tmplErr := loadLoginTemplate()
		if tmplErr != nil {
			c.String(http.StatusInternalServerError, "Error loading template: %v", tmplErr)
			return
		}
		tmpl.Execute(c.Writer, gin.H{"Error": "Invalid credentials", "redirect": redirect})
		return
	}
	// Set session (for demo, use cookie)
	c.SetCookie("session_user", user.Email, 3600, "/", "", false, true)
	if redirect != "" {
		// URL-decode the redirect parameter to restore the original PKCE authorize URL
		decodedRedirect, err := url.QueryUnescape(redirect)
		if err != nil {
			// If decoding fails, use the original redirect
			decodedRedirect = redirect
		}
		c.Redirect(http.StatusFound, decodedRedirect)
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

func Logout(c *gin.Context) {
	// Clear the session cookie
	c.SetCookie("session_user", "", -1, "/", "", false, true)
	redirect := c.Query("redirect")
	if redirect != "" {
		c.Redirect(http.StatusFound, redirect)
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}
