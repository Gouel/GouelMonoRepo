package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/token/auth" || c.Request.URL.Path == "/token/auth/ticket" {
			c.Next()
			return
		}

		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token manquant ou invalide"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Vérifier que le token est signé avec l'algorithme correct
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || token == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userId", claims["userId"])
			c.Set("role", claims["role"])
			c.Set("events", claims["events"])
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Next()
	}
}

func RoleAuthorizationMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé"})
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		for _, role := range requiredRoles {
			if userRole == role {
				c.Next() // Le rôle correspond, continue avec le prochain handler
				return
			}
		}

		// Si c'est une requête avec un event On peut regarder si on a le droit admin

		eventId := c.Param("event_id")
		eventsValue, eventsExists := c.Get("events")

		var events []interface{}
		if eventsExists && eventsValue != nil {
			events = eventsValue.([]interface{})
		}

		if eventId != "" {
			for _, e := range events {
				event := e.(map[string]interface{})
				if event["event_id"] == eventId {
					if event["role"] == "ADMIN" {
						c.Next() // Le rôle correspond, continue avec le prochain handler
						return
					}
				}
			}
		}

		// Si aucun des rôles ne correspond
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé"})
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func EventAuthorizationMiddleware(requiredRights ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, roleExists := c.Get("role")
		eventsValue, eventsExists := c.Get("events")

		if !roleExists || !eventsExists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé", "code": 0x1})
			c.Abort()
			return
		}

		role := roleValue.(string)

		// Accès autorisé pour les superadmins
		if role == "SUPERADMIN" || role == "API" {
			c.Next()
			return
		}

		eventId := c.Param("event_id")
		hasAccess := false

		var events []interface{}
		if eventsExists && eventsValue != nil {
			events = eventsValue.([]interface{})
		}
		for _, e := range events {
			event := e.(map[string]interface{})
			if event["EventId"] == eventId {
				if event["Role"] == "Admin" {
					hasAccess = true
					break
				}

				permissions := event["Permissions"].([]interface{})
				for _, p := range permissions {
					for _, requiredRight := range requiredRights {
						if p == requiredRight {
							hasAccess = true
							break
						}
					}
					if hasAccess {
						break
					}
				}
				if hasAccess {
					break
				}
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé", "code": 0x3})
			c.Abort()
			return
		}

		c.Next()
	}
}

func EventAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, roleExists := c.Get("role")
		eventsValue, eventsExists := c.Get("events")

		if !roleExists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé", "code": 0x1})
			c.Abort()
			return
		}

		role := roleValue.(string)
		var events []interface{}
		if eventsExists && eventsValue != nil {
			events = eventsValue.([]interface{})
		}

		// Accès autorisé pour les rôles API et SUPERADMIN
		if role == "API" || role == "SUPERADMIN" {
			c.Next()
			return
		}

		eventId := c.Param("event_id")
		hasAccess := false

		// Vérifier si l'utilisateur est Administrateur ou Volunteer de l'event
		for _, e := range events {
			event := e.(map[string]interface{})
			if event["EventId"] == eventId {
				if event["Role"] == "Admin" || event["Role"] == "Volunteer" {
					hasAccess = true
					break
				}
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé", "code": 0x3})
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearToken := c.GetHeader("Authorization")
	// Diviser le token à partir du format `Bearer {token}`
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
