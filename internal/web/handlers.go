package web

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"time"
	"uptime-monitor/internal/database/db"
	"uptime-monitor/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

// ShowLoginPage renders the login page.
func (s *Server) ShowLoginPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("internal/web/templates/layout.html", "internal/web/templates/login.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering page: %v", err)
		return
	}

	data := gin.H{
		"title": "Login",
	}

	c.Status(http.StatusOK)
	tmpl.Execute(c.Writer, data)
}

// PostLoginPage handles the login form submission.
func (s *Server) PostLoginPage(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := s.q.GetUserByEmail(context.Background(), email)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.Redirect(http.StatusFound, "/login?error=invalid_credentials")
			return
		}
		c.Redirect(http.StatusFound, "/login?error=database_error")
		return
	}

	match := models.CheckPasswordHash(password, user.PasswordHash)
	if !match {
		c.Redirect(http.StatusFound, "/login?error=invalid_credentials")
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=token_error")
		return
	}

	// Set cookie
	c.SetCookie(CookieName, tokenString, 3600*24*7, "/", "", false, true)
	c.Redirect(http.StatusFound, "/dashboard")
}

// ShowDashboardPage renders the main dashboard page.
func (s *Server) ShowDashboardPage(c *gin.Context) {
	userID := c.GetInt64("userID")

	services, err := s.q.GetServicesForUser(context.Background(), userID)
	if err != nil {
		// Handle error - maybe render a dashboard with an error message
		c.String(http.StatusInternalServerError, "Error fetching services: %v", err)
		return
	}

	tmpl, err := template.ParseFiles("internal/web/templates/layout.html", "internal/web/templates/dashboard.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering page: %v", err)
		return
	}

	data := gin.H{
		"title":    "Dashboard",
		"Services": services,
	}

	c.Status(http.StatusOK)
	tmpl.Execute(c.Writer, data)
}

// Logout handles user logout.
func (s *Server) Logout(c *gin.Context) {
	c.SetCookie(CookieName, "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
