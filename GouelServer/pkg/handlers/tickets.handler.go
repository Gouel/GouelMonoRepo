package handlers

import (
	"net/http"

	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/database"
	"github.com/gin-gonic/gin"
)

// CreateTicketHandler crée un nouveau ticket
func CreateTicketHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	ticketCode := c.Param("ticket_code")

	var requestData struct {
		UserId       string `json:"UserId"`
		WasPurchased *bool  `json:"WasPurchased,omitempty"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	if requestData.WasPurchased == nil {
		requestData.WasPurchased = new(bool)
		*requestData.WasPurchased = true
	}

	ticketId, err := database.CreateTicket(requestData.UserId, eventId, ticketCode, *requestData.WasPurchased)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"TicketId": ticketId})
}

// DeleteTicketHandler supprime un ticket
func DeleteTicketHandler(c *gin.Context) {
	ticketId := c.Param("ticket_id")

	err := database.DeleteTicket(ticketId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket supprimé avec succès"})
}

func ValidateTicketHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	var requestData struct {
		TicketId string `json:"TicketId"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	code, err := database.ValidateTicket(requestData.TicketId, eventId)
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
	eventId := c.Param("event_id")
	var requestData struct {
		TicketId string `json:"TicketId"`
		IsSam    bool   `json:"IsSam"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.SetSAM(requestData.TicketId, eventId, requestData.IsSam)
	if err != nil {
		// Autres erreurs
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket mis à jour"})
}

// GetTicketInfoHandler gère la récupération des informations d'un ticket
func GetTicketInfoHandler(c *gin.Context) {
	ticketId := c.Param("ticket_id")
	eventId := c.Param("event_id")

	ticketInfo, err := database.GetTicketInfo(ticketId, &eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ticketInfo)
}

func GetAllTicketsFromEventHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	tickets, err := database.GetAllTicketsFromEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}
