package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gouel/gouel_serveur/pkg/database"
	"github.com/gouel/gouel_serveur/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

// FindUsersByEmailStartsWithHandler gère la recherche d'utilisateurs par email
func FindUsersByEmailStartsWithHandler(c *gin.Context) {
	email := c.Param("email")
	users, err := database.FindUsersByEmailStartsWith(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByIDHandler gère la récupération d'un utilisateur par son ID
func GetUserByIDHandler(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := database.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUserHandler gère la création d'un nouvel utilisateur
func CreateUserHandler(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := database.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func UpdateUserHandler(c *gin.Context) {
	userID := c.Param("user_id")

	var updateData bson.M
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mettre à jour l'utilisateur dans la base de données
	err := database.UpdateUser(userID, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur mis à jour"})
}

func AddUserTransactionHandler(c *gin.Context) {
	userID := c.Param("user_id")

	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddUserTransaction(userID, transaction)
	if err != nil {
		// Gérer les erreurs spécifiques retournées par AddUserTransaction
		if err.Error() == "solde insuffisant" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Solde insuffisant"})
			return
		} else if err.Error() == "type de transaction non reconnu" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Type de transaction non reconnu"})
			return
		}
		// Autres erreurs non spécifiques
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction ajoutée avec succès"})
}

func UserPayHandler(c *gin.Context) {
	userID := c.Param("user_id")
	event_id := c.Param("event_id")

	var purchaseItems models.PurchaseItems
	if err := c.ShouldBindJSON(&purchaseItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	// Logique pour traiter l'achat
	totalCost, err := processPurchase(event_id, userID, purchaseItems)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Générer une transaction de débit
	err = addDebitTransaction(userID, totalCost, purchaseItems)
	if err != nil {
		if err.Error() == "solde insuffisant" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Achat effectué avec succès", "total_cost": totalCost})
}

func processPurchase(event_id string, userID string, items models.PurchaseItems) (float64, error) {
	var totalCost float64
	user, err := database.GetUserByID(userID)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		product, err := database.GetEventProductsByCode(event_id, item.ProductCode)
		if err != nil {
			return 0, err
		}

		if hasEnded(product.End_Of_Sale) {
			return 0, fmt.Errorf("le produit %s n'est plus en vente", item.ProductCode)
		}

		if product.Is_Adult && !isUserAdult(user["DOB"].(string)) {
			return 0, fmt.Errorf("l'utilisateur n'a pas l'âge requis pour acheter le produit %s", item.ProductCode)
		}

		totalCost += product.Price * float64(item.Amount)
	}
	return totalCost, nil
}

func isUserAdult(dob string) bool {
	// Convertir la chaîne de date de naissance en time.Time
	dobTime, err := time.Parse("2006-01-02", dob)
	if err != nil {
		fmt.Println("Erreur lors de la conversion de la date de naissance :", err)
		return false // Retourner false si la date de naissance n'est pas valide
	}

	// Calculer l'âge en comparant avec la date actuelle
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)
	return dobTime.Before(eighteenYearsAgo)
}

func hasEnded(endOfSale string) bool {
	// Convertir la chaîne de date en time.Time
	endOfSaleTime, err := time.Parse("2006-01-02T15:04", endOfSale)
	if err != nil {
		fmt.Println("Erreur lors de la conversion de la date :", err)
		return true // Considérer comme terminé si la date n'est pas valide
	}

	// Comparer avec la date actuelle
	return time.Now().After(endOfSaleTime)
}

func addDebitTransaction(userID string, amount float64, cart models.PurchaseItems) error {
	// Créer une transaction de débit
	debitTransaction := models.Transaction{
		Type:   "debit",
		Amount: amount,
		Cart:   cart,
		Date:   time.Now().Format("2006-01-02T15:04:05"),
	}

	// Utiliser AddUserTransaction pour mettre à jour le solde de l'utilisateur et enregistrer la transaction
	err := database.AddUserTransaction(userID, debitTransaction)
	if err != nil {
		return err // Gérer les erreurs, par exemple, solde insuffisant
	}

	return nil
}
