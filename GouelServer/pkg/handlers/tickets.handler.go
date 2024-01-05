package handlers

import (
	// ... autres imports

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gouel/gouel_serveur/pkg/database"
)

// CreateTicketHandler crée un nouveau ticket
func CreateTicketHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	ticketCode := c.Param("ticket_code")

	var requestData struct {
		UserID    string `json:"user_id"`
		IsSpecial *bool  `json:"is_special,omitempty"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	if requestData.IsSpecial == nil {
		requestData.IsSpecial = new(bool)
		*requestData.IsSpecial = false
	}

	ticketID, err := database.CreateTicket(requestData.UserID, eventID, ticketCode, *requestData.IsSpecial)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ticket_id": ticketID})
}

// DeleteTicketHandler supprime un ticket
func DeleteTicketHandler(c *gin.Context) {
	ticketID := c.Param("ticket_id")

	err := database.DeleteTicket(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket supprimé avec succès"})
}

func ValidateTicketHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	var requestData struct {
		TicketID string `json:"ticket_id"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	code, err := database.ValidateTicket(requestData.TicketID, eventID)
	if err != nil {
		// Vérifier si l'erreur est due à un ticket déjà validé
		if code == 2 {
			c.JSON(http.StatusConflict, gin.H{"error": "Le ticket a déjà été validé"})
			return
		}
		// Autres erreurs
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket validé avec succès"})
}

func SetSamHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	var requestData struct {
		TicketID string `json:"ticket_id"`
		IsSam    bool   `json:"is_sam"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.SetSAM(requestData.TicketID, eventID, requestData.IsSam)
	if err != nil {
		// Autres erreurs
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket mis à jour"})
}

// GetTicketInfoHandler gère la récupération des informations d'un ticket
func GetTicketInfoHandler(c *gin.Context) {
	ticketID := c.Param("ticket_id")
	eventID := c.Param("event_id")

	ticketInfo, err := database.GetTicketInfo(ticketID, &eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ticketInfo)
}

func GetAllTicketsFromEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	tickets, err := database.GetAllTicketsFromEvent(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}
