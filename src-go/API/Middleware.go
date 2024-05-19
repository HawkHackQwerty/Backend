package Endpoints

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/ratelimit"
	"net/http"
	"time"
)

var limit = ratelimit.New(100) // const 100 requests per second allowed

func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(c *gin.Context) {
		now := limit.Take()
		fmt.Printf("LeakBucket Middleware: Time since last request: %v\n", now.Sub(prev))
		prev = now
		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			fmt.Println("Auth Middleware: Authorization token not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			return
		}

		tokenString = tokenString[7:] // Remove "Bearer " from the start
		fmt.Printf("Auth Middleware: Token after stripping Bearer: %s\n", tokenString)

		// Hardcoded secret key
		secret := []byte("C0YrhqJLfRMfjsOL75ahnObNCc5D4UZ/G6NdNTse1LH2OfX/3qnab3zWUegeKhr5oWlRCCNe5kONfMaYCOoMDQ==")

		// Parse and verify the JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		if err != nil {
			fmt.Printf("Auth Middleware: Error parsing token: %s\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			return
		}

		if !token.Valid {
			fmt.Println("Auth Middleware: Token is invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}

// Middleware to extract user ID from the X-User-ID header
func userIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			fmt.Println("UserID Middleware: User ID not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not provided"})
			return
		}

		fmt.Printf("UserID Middleware: User ID from header: %s\n", userID)
		c.Set("userID", userID)
		c.Next()
	}
}
