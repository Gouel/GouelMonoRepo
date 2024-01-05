package token

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gouel/gouel_serveur/pkg/config"
	"github.com/gouel/gouel_serveur/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
)

func AuthRoute(c *gin.Context) {
	// Charger la configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne du serveur"})
		return
	}

	// Récupération de l'email et du mot de passe depuis la requête
	var loginInfo struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	// Vérifier l'utilisateur
	user, err := database.VerifyUser(loginInfo.Email, loginInfo.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Couple email/password invalide"})
		return
	}

	// Générer le token JWT
	token, err := createToken(user.ID.Hex(), user.Role, cfg.JWTSecretKey, int64(cfg.JWTExpiration))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
func AuthRouteTicket(c *gin.Context) {
	// Charger la configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne du serveur"})
		return
	}

	// Récupération de l'user_id depuis la requête
	var loginInfo struct {
		TicketID string `json:"ticket_id"` // Modifié pour commencer par une majuscule
	}
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	// Vérifier le ticket
	ticketInfo, err := database.GetTicketInfo(loginInfo.TicketID, nil)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ticket invalid"})
		return
	}

	userID := ticketInfo["user"].(bson.M)["user_id"].(string)
	// Vérifier l'user
	user, err := database.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user invalid"})
		return
	}

	// Générer le token JWT
	token, err := createToken(userID, user["role"].(string), cfg.JWTSecretKey, int64(cfg.JWTExpiration))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func RefreshRoute(c *gin.Context) {
	// Charger la configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la configuration: %v", err)
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	token, err := createToken(userID.(string), role.(string), cfg.JWTSecretKey, int64(cfg.JWTExpiration))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func createToken(userID, role, secretKey string, expirationMinutes int64) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)

	eventsRoles, err := database.GetUserEventsRoles(userID)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     expirationTime.Unix(),
		"events":  eventsRoles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ViewTokenRoute(c *gin.Context) {
	tokenString := extractToken(c)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant ou invalide"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
		}
		secretKey := os.Getenv("JWT_SECRET_KEY") // Assurez-vous que cette clé est définie dans votre .env
		return []byte(secretKey), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.JSON(http.StatusOK, claims)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide ou expiré"})
	}
}

func extractToken(c *gin.Context) string {
	bearToken := c.GetHeader("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
