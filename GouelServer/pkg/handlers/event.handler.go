package handlers

import (
	"net/http"

	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/config"
	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/database"
	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GET

func GetSMTPConfiguration(c *gin.Context) {
	// Charger la configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne du serveur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Email":         cfg.Email,
		"EmailPassword": cfg.EmailPassword,
		"SMTPServer":    cfg.SMTPServer,
		"SMTPPort":      cfg.SMTPPort,
		"SMTPSSL":       cfg.SMTPUseSSL,
	})
}

func GetAccessibleEventsHandler(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identifiant utilisateur non trouvé"})
		return
	}

	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Rôle utilisateur non trouvé"})
		return
	}

	events, err := database.GetAccessibleEvents(userId.(string), userRole.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		println("DEBUG", err.Error())
		return
	}

	c.JSON(http.StatusOK, events)
}

func GetEventByIdHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	event, err := database.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}
func GetSimpleEventHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	event, err := database.GetSimpleEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}
func GetEventProductsHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	products, err := database.GetEventProducts(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
func GetEventProductsCodeHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	product_code := c.Param("product_code")
	products, err := database.GetEventProductsByCode(eventId, product_code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
func GetEventAdminsHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	admins, err := database.GetEventAdmins(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, admins)
}
func GetEventVolunteersHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	volunteers, err := database.GetEventVolunteers(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, volunteers)
}

func GetEventLockersHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	lockers, err := database.GetEventLockers(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lockers)
}

func GetEventTicketsHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	tickets, err := database.GetEventTickets(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

// POST

func AddEventHandler(c *gin.Context) {
	var newEvent models.Event
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	eventId, err := database.AddEvent(newEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event_id": eventId})
}

func AddTicketToEventHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	var ticketData models.EventTicket
	if err := c.ShouldBindJSON(&ticketData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddEventTicketToEvent(eventId, ticketData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket ajouté avec succès"})
}
func AddLockerToEventHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	var lockerData models.Locker
	if err := c.ShouldBindJSON(&lockerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddLockerToEvent(eventId, lockerData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Casier ajouté avec succès"})
}

func AddRangeOfLockersHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	var lockerRangeData struct {
		Start  int    `json:"start"`
		End    int    `json:"end"`
		Prefix string `json:"prefix"`
		Suffix string `json:"suffix"`
	}
	if err := c.ShouldBindJSON(&lockerRangeData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	lockers, err := database.AddRangeOfLockersToEvent(eventId, lockerRangeData.Start, lockerRangeData.End, lockerRangeData.Prefix, lockerRangeData.Suffix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lockers)
}

func AddVolunteerToEventHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	var volunteerData models.Volunteer
	if err := c.ShouldBindJSON(&volunteerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddVolunteerToEvent(eventId, volunteerData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bénévole ajouté avec succès"})
}

func AddProductToEventHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	var productData models.Product
	if err := c.ShouldBindJSON(&productData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddProductToEvent(eventId, productData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produit ajouté avec succès"})
}

// PUT

func UpdateEventHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	var eventData bson.M
	if err := c.ShouldBindJSON(&eventData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.UpdateEvent(eventId, eventData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Événement mis à jour avec succès"})
}

func UpdateLockerHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	var lockerData models.Locker
	if err := c.ShouldBindJSON(&lockerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.UpdateLocker(eventId, lockerData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Casier mis à jour avec succès"})
}

func UpdateProductHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	product_code := c.Param("product_code")

	var productData bson.M
	if err := c.ShouldBindJSON(&productData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.UpdateProduct(eventId, product_code, productData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produit mis à jour avec succès"})
}

func UpdateVolunteerHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	var volunteer models.Volunteer
	if err := c.ShouldBindJSON(&volunteer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}
	err := database.UpdateVolunteer(eventId, volunteer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Permissions du bénévole mises à jour avec succès"})
}

func UpdateTicketHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	ticketCode := c.Param("ticket_code")

	var ticketData bson.M
	if err := c.ShouldBindJSON(&ticketData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.UpdateEventTicket(eventId, ticketCode, ticketData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket mis à jour avec succès"})
}

// DELETE

func DeleteEventHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	err := database.DeleteEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Événement supprimé avec succès"})
}
func DeleteEventTicketHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	ticketCode := c.Param("ticket_code")

	err := database.DeleteEventTicket(eventId, ticketCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket supprimé avec succès"})
}
func DeleteLockerHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	lockerCode := c.Param("locker_code")

	err := database.DeleteLocker(eventId, lockerCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Casier supprimé avec succès"})
}
func DeleteAllLockersHandler(c *gin.Context) {
	eventId := c.Param("event_id")

	err := database.DeleteAllLockers(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tous les casiers ont été supprimés avec succès"})
}
func DeleteVolunteerHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	userId := c.Param("user_id")

	err := database.DeleteVolunteer(eventId, userId)
	if err != nil {
		// Intercepter l'erreur spécifique et renvoyer un code d'erreur approprié
		if err.Error() == "impossible de supprimer un bénévole qui est également administrateur" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bénévole supprimé avec succès"})
}

func DeleteProductHandler(c *gin.Context) {
	eventId := c.Param("event_id")
	productCode := c.Param("product_code")

	err := database.DeleteProduct(eventId, productCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produit supprimé avec succès"})
}
