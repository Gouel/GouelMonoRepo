package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gouel/gouel_serveur/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// <=== Lecture de la BDD ===>

func GetAllEventsIDs() ([]string, error) {
	var ids []string
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var event bson.M
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		id := event["_id"].(primitive.ObjectID).Hex()
		ids = append(ids, id)
	}

	return ids, nil
}

func GetAccessibleEvents(userID, userRole string) ([]models.Event, error) {
	var events []models.Event
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	if userRole != "SUPERADMIN" && userRole != "API" {
		// Filtre pour les utilisateurs qui ne sont ni SUPERADMIN ni API
		filter = bson.M{
			"$or": []bson.M{
				{"public": true},
				{"admins": userID},
			},
		}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var event models.Event
		if err := cursor.Decode(&event); err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(eventID string) (bson.M, error) {
	var event bson.M
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func GetSimpleEvent(eventID string) (bson.M, error) {
	var event bson.M
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	err := collection.FindOne(ctx, bson.M{"_id": objID}, options.FindOne().SetProjection(bson.M{"title": 1, "public": 1, "description": 1, "location": 1, "contact": 1, "image_url": 1})).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func GetEventProducts(eventID string) (interface{}, error) {
	var event bson.M
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return event["products"], nil
}

func GetEventProductsByCode(eventID, productCode string) (*models.Product, error) {
	var event models.Event
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return nil, err
	}

	// Récupérer l'événement spécifique
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return nil, err
	}

	// Rechercher le produit par son code dans la liste des produits de l'événement
	for _, product := range event.Products {
		if product.Code == productCode {
			return &product, nil
		}
	}

	return nil, fmt.Errorf("produit avec le code %s non trouvé dans l'événement", productCode)
}

func GetEventAdmins(eventID string) (interface{}, error) {
	var event bson.M
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event["admins"], nil
}

func GetEventVolunteers(eventID string) (interface{}, error) {
	var event bson.M
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event["volunteers"], nil
}

func GetEventLockers(eventID string) (interface{}, error) {
	var event bson.M
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event["lockers"], nil
}

func GetEventTickets(eventID string) (interface{}, error) {
	var event bson.M
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event["tickets"], nil
}

// pour token

func GetUserEventsRoles(userID string) ([]models.EventRole, error) {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"volunteers.user": userID},
			{"admins": userID},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventsRoles []models.EventRole
	for cursor.Next(ctx) {
		var event models.Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}

		var role string
		var permissions []string

		if isUserInAdmins(event.Admins, userID) {
			role = "Admin"
			permissions = make([]string, 0)
		} else {
			role = "Volunteer"
			volunteer, found := getVolunteerFromEvent(event.Volunteers, userID)
			if found {
				permissions = volunteer.Permissions
			}
		}

		eventsRoles = append(eventsRoles, models.EventRole{
			EventID: event.ID.Hex(),
			Role:    role,
			Rights:  permissions,
		})
	}

	return eventsRoles, nil
}

func isUserInAdmins(admins []string, userID string) bool {
	for _, admin := range admins {
		if admin == userID {
			return true
		}
	}
	return false
}

func getVolunteerFromEvent(volunteers []models.Volunteer, userID string) (models.Volunteer, bool) {
	for _, volunteer := range volunteers {
		if volunteer.User_ID == userID {
			return volunteer, true
		}
	}
	return models.Volunteer{}, false
}

// <=== Ajout dans la BDD ===>

func AddEvent(event models.Event) (string, error) {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialiser les champs à des valeurs par défaut
	event.Admins = []string{}
	event.Volunteers = make([]models.Volunteer, 0)
	event.Lockers = make([]models.Locker, 0)
	event.Products = make([]models.Product, 0)
	event.Tickets = make([]models.EventTicket, 0)

	result, err := collection.InsertOne(ctx, event)
	id := ""
	if err == nil {
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			id = oid.Hex()
		}
	}
	return id, err
}

func AddAdminToEvent(eventID string, admin string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	update := bson.M{"$push": bson.M{"admins": admin}}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func AddTicketToEvent(eventID string, ticket bson.M) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)

	code := uuid.New().String()
	ticket["code"] = code

	update := bson.M{"$push": bson.M{"tickets": ticket}}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func AddVolunteerToEvent(eventID string, volunteer bson.M) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	update := bson.M{"$push": bson.M{"volunteers": volunteer}}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func AddLockerToEvent(eventID string, locker bson.M) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	update := bson.M{"$push": bson.M{"lockers": locker}}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func AddRangeOfLockersToEvent(eventID string, start, end int, prefix, suffix string) ([]models.Locker, error) {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Déterminer le nombre de chiffres nécessaire pour le formatage
	digitCount := len(fmt.Sprintf("%d", end))

	// Construire la liste des casiers
	var lockers []models.Locker
	for i := start; i <= end; i++ {
		// Générer le code avec le nombre de chiffres constants
		lockerCode := fmt.Sprintf("%s%0*d%s", prefix, digitCount, i, suffix)
		lockers = append(lockers, models.Locker{Code: lockerCode, User: ""})
	}

	// Préparer la mise à jour pour ajouter les casiers
	update := bson.M{
		"$push": bson.M{"lockers": bson.M{"$each": lockers}},
	}

	objID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return nil, err
	}

	// Exécuter la mise à jour
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return nil, err
	}

	return lockers, nil
}

func AddProductToEvent(eventID string, product bson.M) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)

	code := uuid.New().String()
	product["code"] = code

	update := bson.M{"$push": bson.M{"products": product}}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// <=== Modif dans la BDD ===>

func UpdateEvent(eventID string, updateData bson.M) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	return err
}

func UpdateLocker(eventID string, lockerCode string, userID *string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)

	var update bson.M
	if userID != nil {
		update = bson.M{"$set": bson.M{"lockers.$.user": *userID}}
	} else {
		update = bson.M{"$set": bson.M{"lockers.$.user": nil}}
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID, "lockers.code": lockerCode}, update)
	return err
}

func UpdateProduct(eventID, productCode string, updateData bson.M) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)

	// Créer un nouveau bson.M pour la mise à jour avec l'opérateur $.
	setUpdateData := bson.M{}
	for key, value := range updateData {
		setUpdateData["products.$."+key] = value
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID, "products.code": productCode}, bson.M{"$set": setUpdateData})
	return err
}

func UpdateVolunteerPermissions(eventID string, userID string, permissions []string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID, "volunteers.user_id": userID}, bson.M{"$set": bson.M{"volunteers.$.permissions": permissions}})
	return err
}

func UpdateTicket(eventID, ticketCode string, updateData bson.M) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)

	// Créer un nouveau bson.M pour la mise à jour avec l'opérateur $.
	setUpdateData := bson.M{}
	for key, value := range updateData {
		setUpdateData["tickets.$."+key] = value
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID, "tickets.code": ticketCode}, bson.M{"$set": setUpdateData})
	return err
}

// <=== Suppression dans la BDD ===>

func DeleteEvent(eventID string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func DeleteEventTicket(eventID, ticketCode string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$pull": bson.M{"tickets": bson.M{"code": ticketCode}}})
	return err
}

func DeleteLocker(eventID, lockerCode string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$pull": bson.M{"lockers": bson.M{"code": lockerCode}}})
	return err
}

func DeleteAllLockers(eventID string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convertir l'ID de l'événement en ObjectID
	objID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return err
	}

	// Mise à jour pour supprimer tous les casiers
	update := bson.M{
		"$set": bson.M{"lockers": []models.Locker{}}, // Définir les casiers à une liste vide
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	return nil
}

func DeleteVolunteer(eventID, userID string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return err
	}

	// Vérifier si l'utilisateur est un administrateur de l'événement
	var event models.Event
	if err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event); err != nil {
		return err
	}

	for _, adminID := range event.Admins {
		if adminID == userID {
			return fmt.Errorf("impossible de supprimer un bénévole qui est également administrateur")
		}
	}

	// Supprimer le bénévole si ce n'est pas un administrateur
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$pull": bson.M{"volunteers": bson.M{"user_id": userID}}})
	return err
}

func DeleteProduct(eventID, productCode string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$pull": bson.M{"products": bson.M{"code": productCode}}})
	return err
}

func DeleteAdmin(eventID, adminID string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(eventID)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$pull": bson.M{"admins": adminID}})
	return err
}
