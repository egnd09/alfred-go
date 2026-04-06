package api

import (
	"context"
	"net/http"
	"time"

	"github.com/egnd09/alfred-go/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type LoginRequest struct {
	Code string `json:"code"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

type GitHubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

// Login handles GitHub OAuth login
func Login(cfg *config.Config, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if req.Code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
			return
		}

		// Exchange code for access token
		token, err := exchangeGitHubCode(cfg, req.Code)
		if err != nil {
			logger.Error("failed to exchange code", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "oauth failed"})
			return
		}

		// Get GitHub user info
		user, err := getGitHubUser(token)
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
			return
		}

		// Generate JWT token
		jwtToken, err := generateJWT(user, cfg.JWTSecret)
		if err != nil {
			logger.Error("failed to generate jwt", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token: jwtToken,
			User:  user.Login,
		})
	}
}

func exchangeGitHubCode(cfg *config.Config, code string) (string, error) {
	// TODO: Implement actual GitHub OAuth exchange
	// This is a placeholder - implement actual GitHub OAuth flow
	return "mock_token", nil
}

func getGitHubUser(token string) (*GitHubUser, error) {
	// TODO: Implement actual GitHub user API call
	// This is a placeholder - implement actual API call
	return &GitHubUser{
		ID:    1,
		Login: "testuser",
	}, nil
}

func generateJWT(user *GitHubUser, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"login":   user.Login,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// AuthMiddleware validates JWT token
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := ValidateJWT(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Store claims in context
		c.Set("claims", claims)
		c.Next()
	}
}
