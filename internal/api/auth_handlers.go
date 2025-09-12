package api

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"
	"uptime-monitor/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

// userInput es una estructura para recibir los datos del cliente.
type userInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// registerUser maneja el registro de nuevos usuarios.
func (s *Server) registerUser(c *gin.Context) {
	var input userInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user := models.User{
		Email: input.Email,
	}

	// Hashear la contraseña
	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Guardar en la base de datos
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := s.db.QueryRow(context.Background(), query, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		// Manejar error de email duplicado
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": user.ID})
}

// loginUser maneja la autenticación y devuelve un token JWT.
func (s *Server) loginUser(c *gin.Context) {
	var input userInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var user models.User
	query := `SELECT id, password_hash FROM users WHERE email = $1`
	err := s.db.QueryRow(context.Background(), query, input.Email).Scan(&user.ID, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Verificar la contraseña
	match, err := user.CheckPassword(input.Password)
	if err != nil || !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Crear el token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,                                   // El ID del usuario bajo la clave 'user_id'
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // Expira en 7 días
	})

	// Firmar el token con el secreto
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
