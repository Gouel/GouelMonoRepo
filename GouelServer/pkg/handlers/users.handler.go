package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/database"
	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/models"
	"github.com/gin-gonic/gin"
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
func GetUserByIdHandler(c *gin.Context) {
	userId := c.Param("user_id")
	fmt.Println(userId)
	user, err := database.GetUserById(userId)
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
	userId, err := database.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"UserId": userId})
}

func UpdateUserHandler(c *gin.Context) {
	userId := c.Param("user_id")

	var updateData bson.M
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mettre à jour l'utilisateur dans la base de données
	err := database.UpdateUser(userId, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur mis à jour"})
}

func AddUserTransactionHandler(c *gin.Context) {
	userId := c.Param("user_id")

	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddUserTransaction(userId, transaction)
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
	ticketId := c.Param("ticket_id")
	eventId := c.Param("event_id")

	var purchaseItems []models.PurchaseProduct
	if err := c.ShouldBindJSON(&purchaseItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides", "code": 0x01})
		return
	}

	ticket, err := database.GetTicketInfo(ticketId, &eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket invalide", "code": 0x02})
		return
	}

	// Logique pour traiter l'achat
	totalCost, err2 := processPurchase(eventId, ticket, purchaseItems)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2, "code": 0x03})
		return
	}

	// Générer une transaction de débit
	err = addDebitTransaction(ticket.UserId, totalCost, purchaseItems)
	if err != nil {
		if err.Error() == "solde insuffisant" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": 0x04})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": 0xFF})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Achat effectué avec succès", "total_cost": totalCost, "code": 0x00})
}

func processPurchase(eventId string, ticket *models.Ticket, items []models.PurchaseProduct) (float32, *models.GouelError) {
	var totalCost float32

	for _, item := range items {
		product, err := database.GetEventProductsByCode(eventId, item.ProductCode)
		if err != nil {
			return 0, &models.GouelError{
				Message: "Produit invalide",
				Code:    0x00,
				Data:    item.ProductCode,
			}
		}

		if hasEnded(product.EndOfSale) {
			return 0, &models.GouelError{
				Message: "La vente de ce produit est terminée",
				Code:    0x01,
				Data:    product,
			}
		}

		if product.HasAlcohol && !isUserAdult(*ticket) {
			return 0, &models.GouelError{
				Message: "L'utilisateur n'a pas le droit d'acheter ce produit",
				Code:    0x02,
				Data:    product,
			}
		}

		totalCost += product.Price * float32(item.Amount)
	}
	return totalCost, nil
}

func isUserAdult(ticket models.Ticket) bool {

	if ticket.IsSam {
		return false
	}

	// Convertir la chaîne de date de naissance en time.Time
	dobTime, err := time.Parse("2006-01-02", ticket.User.DOB)
	if err != nil {
		fmt.Println("Erreur lors de la conversion de la date de naissance :", err)
		return false // Retourner false si la date de naissance n'est pas valide
	}

	// Calculer l'âge en comparant avec la date actuelle
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)
	return dobTime.Before(eighteenYearsAgo)
}

func hasEnded(endOfSale string) bool {

	if endOfSale == "" {
		return false
	}

	// Convertir la chaîne de date en time.Time
	endOfSaleTime, err := time.Parse("2006-01-02T15:04", endOfSale)
	if err != nil {
		fmt.Println("Erreur lors de la conversion de la date :", err)
		return true // Considérer comme terminé si la date n'est pas valide
	}

	// Comparer avec la date actuelle
	return time.Now().After(endOfSaleTime)
}

func addDebitTransaction(userId string, amount float32, cart []models.PurchaseProduct) error {
	// Créer une transaction de débit
	debitTransaction := models.Transaction{
		Type:   "debit",
		Amount: amount,
		Cart:   cart,
		Date:   time.Now().Format("2006-01-02T15:04:05"),
	}

	// Utiliser AddUserTransaction pour mettre à jour le solde de l'utilisateur et enregistrer la transaction
	err := database.AddUserTransaction(userId, debitTransaction)
	if err != nil {
		return err // Gérer les erreurs, par exemple, solde insuffisant
	}

	return nil
}
