package handlers

import (
	"net/http"
	"strconv"

	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/database"
	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/models"
	"github.com/gin-gonic/gin"
)

// CreateTicketHandler crée un nouveau ticket
func CreateTicketHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	ticketCode := c.Param("ticket_code")
	var TicketRequestData models.TicketRequestData
	if err := c.ShouldBindJSON(&TicketRequestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	database.LogAction(c, TicketRequestData)

	if TicketRequestData.WasPurchased == nil {
		TicketRequestData.WasPurchased = new(bool)
		*TicketRequestData.WasPurchased = true
	}

	if TicketRequestData.PurchasedOnline == nil {
		TicketRequestData.PurchasedOnline = new(bool)
		*TicketRequestData.PurchasedOnline = false
	}

	ticketId, err := database.CreateTicket(TicketRequestData.UserId, eventId, ticketCode, TicketRequestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"TicketId": ticketId})
}

// DeleteTicketHandler supprime un ticket
func DeleteTicketHandler(c *gin.Context) {
	ticketId := c.Param("ticket_id")
	database.LogAction(c, nil)

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
	database.LogAction(c, requestData)

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

	database.LogAction(c, requestData)

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

func GetPaginatedTicketsFromEventHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	page := c.Param("page")

	// Vérifier si la page est un nombre
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page is not a number"})
		return
	}

	tickets, err := database.GetPaginatedTicketsFromEvent(eventId, int64(pageInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func ReturnEcoCupHandler(c *gin.Context) {
	var requestData struct {
		TicketId string `json:"TicketId"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	database.LogAction(c, requestData)

	err := database.ReturnEcoCup(requestData.TicketId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "EcoCup retournée avec succès"})
}
