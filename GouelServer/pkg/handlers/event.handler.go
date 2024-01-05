package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gouel/gouel_serveur/pkg/database"
	"github.com/gouel/gouel_serveur/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

// GET

func GetAllEventsIDsHandler(c *gin.Context) {
	// Logique pour obtenir tous les IDs d'événements
	eventIDs, err := database.GetAllEventsIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, eventIDs)
}

func GetAccessibleEventsHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identifiant utilisateur non trouvé"})
		return
	}

	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Rôle utilisateur non trouvé"})
		return
	}

	events, err := database.GetAccessibleEvents(userID.(string), userRole.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func GetEventByIDHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	event, err := database.GetEventByID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}
func GetSimpleEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	event, err := database.GetSimpleEvent(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}
func GetEventProductsHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	products, err := database.GetEventProducts(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
func GetEventProductsCodeHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	product_code := c.Param("product_code")
	products, err := database.GetEventProductsByCode(eventID, product_code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
func GetEventAdminsHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	admins, err := database.GetEventAdmins(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, admins)
}
func GetEventVolunteersHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	volunteers, err := database.GetEventVolunteers(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, volunteers)
}

func GetEventLockersHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	lockers, err := database.GetEventLockers(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lockers)
}

func GetEventTicketsHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	tickets, err := database.GetEventTickets(eventID)
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

	eventID, err := database.AddEvent(newEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event_id": eventID})
}
func AddAdminToEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	var adminData struct {
		AdminID string `json:"admin_id"`
	}
	if err := c.ShouldBindJSON(&adminData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddAdminToEvent(eventID, adminData.AdminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admin ajouté avec succès"})
}
func AddTicketToEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	var ticketData models.EventTicket
	if err := c.ShouldBindJSON(&ticketData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddTicketToEvent(eventID, bson.M{
		"title": ticketData.Title,
		"price": ticketData.Price,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket ajouté avec succès"})
}
func AddLockerToEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	var lockerData models.Locker
	if err := c.ShouldBindJSON(&lockerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddLockerToEvent(eventID, bson.M{
		"code": lockerData.Code,
		"user": lockerData.User,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Casier ajouté avec succès"})
}

func AddRangeOfLockersHandler(c *gin.Context) {
	eventID := c.Param("event_id")

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

	lockers, err := database.AddRangeOfLockersToEvent(eventID, lockerRangeData.Start, lockerRangeData.End, lockerRangeData.Prefix, lockerRangeData.Suffix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lockers)
}

func AddVolunteerToEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	var volunteerData models.Volunteer
	if err := c.ShouldBindJSON(&volunteerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddVolunteerToEvent(eventID, bson.M{
		"user_id":     volunteerData.User_ID,
		"permissions": volunteerData.Permissions,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bénévole ajouté avec succès"})
}

func AddProductToEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	var productData models.Product
	if err := c.ShouldBindJSON(&productData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.AddProductToEvent(eventID, bson.M{
		"icon":        productData.Icon,
		"end_of_sale": productData.End_Of_Sale,
		"label":       productData.Label,
		"price":       productData.Price,
		"adult":       productData.Is_Adult,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produit ajouté avec succès"})
}

// PUT

func UpdateEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	var eventData bson.M
	if err := c.ShouldBindJSON(&eventData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.UpdateEvent(eventID, eventData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Événement mis à jour avec succès"})
}

func UpdateLockerHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	var lockerData models.Locker
	if err := c.ShouldBindJSON(&lockerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.UpdateLocker(eventID, lockerData.Code, &lockerData.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Casier mis à jour avec succès"})
}

func UpdateProductHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	product_code := c.Param("product_code")

	var productData bson.M
	if err := c.ShouldBindJSON(&productData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.UpdateProduct(eventID, product_code, productData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produit mis à jour avec succès"})
}

func UpdateVolunteerPermissionsHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	userID := c.Param("user_id")

	var permissionsData struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&permissionsData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}
	err := database.UpdateVolunteerPermissions(eventID, userID, permissionsData.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Permissions du bénévole mises à jour avec succès"})
}

func UpdateTicketHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	ticketCode := c.Param("ticket_code")

	var ticketData bson.M
	if err := c.ShouldBindJSON(&ticketData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'entrée invalides"})
		return
	}

	err := database.UpdateTicket(eventID, ticketCode, ticketData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket mis à jour avec succès"})
}

// DELETE

func DeleteEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	err := database.DeleteEvent(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Événement supprimé avec succès"})
}
func DeleteEventTicketHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	ticketCode := c.Param("ticket_code")

	err := database.DeleteEventTicket(eventID, ticketCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket supprimé avec succès"})
}
func DeleteLockerHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	lockerCode := c.Param("locker_code")

	err := database.DeleteLocker(eventID, lockerCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Casier supprimé avec succès"})
}
func DeleteAllLockersHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	err := database.DeleteAllLockers(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tous les casiers ont été supprimés avec succès"})
}
func DeleteVolunteerHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	userID := c.Param("user_id")

	err := database.DeleteVolunteer(eventID, userID)
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
	eventID := c.Param("event_id")
	productCode := c.Param("product_code")

	err := database.DeleteProduct(eventID, productCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produit supprimé avec succès"})
}
func DeleteAdminHandler(c *gin.Context) {
	eventID := c.Param("event_id")
	adminID := c.Param("admin_id")

	err := database.DeleteAdmin(eventID, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Administrateur supprimé avec succès"})
}
