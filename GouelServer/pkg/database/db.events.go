package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// <=== Lecture de la BDD ===>

func GetAllEventsIds() ([]string, error) {
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

func GetAccessibleEvents(userId, userRole string) ([]models.Event, error) {
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
				{"volunteers": bson.M{"$elemMatch": bson.M{"userId": userId, "isAdmin": true}}},
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

func GetEventById(eventId string) (bson.M, error) {
	var event bson.M
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func GetSimpleEvent(eventId string) (bson.M, error) {
	var event bson.M
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	err := collection.FindOne(ctx, bson.M{"_id": objId}, options.FindOne().SetProjection(bson.M{"title": 1, "public": 1, "description": 1, "location": 1, "contact": 1})).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func GetEventProducts(eventId string) (interface{}, error) {
	var event bson.M
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return event["products"], nil
}

func GetEventProductsByCode(eventId, productCode string) (*models.Product, error) {
	var event models.Event
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		return nil, err
	}

	// Récupérer l'événement spécifique
	err = collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&event)
	if err != nil {
		return nil, err
	}

	// Rechercher le produit par son code dans la liste des produits de l'événement
	for _, product := range event.Products {
		if product.ProductCode == productCode {
			return &product, nil
		}
	}

	return nil, fmt.Errorf("produit avec le code %s non trouvé dans l'événement", productCode)
}

func GetEventAdmins(eventId string) (interface{}, error) {
	var event models.Event
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&event)
	if err != nil {
		return nil, err
	}

	admins := []models.Volunteer{}

	for _, volunteer := range event.Volunteers {
		if volunteer.IsAdmin {
			admins = append(admins, volunteer)
		}
	}

	return admins, nil
}

func GetEventVolunteers(eventId string) (interface{}, error) {
	var event models.Event
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event.Volunteers, nil
}

func GetEventLockers(eventId string) (interface{}, error) {
	var event models.Event
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event.Lockers, nil
}

func GetEventTickets(eventId string) (interface{}, error) {
	var event models.Event
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event.EventTickets, nil
}

// pour token

func GetUserEventsRoles(userId string) ([]models.EventRole, error) {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"volunteers": bson.M{"$elemMatch": bson.M{"userId": userId}},
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

		volunteer, found := getVolunteerFromEvent(event.Volunteers, userId)
		if found {
			permissions = volunteer.Permissions
			role = "Volunteer"
			if volunteer.IsAdmin {
				role = "Admin"
			}
		} else {
			role = "User"
		}

		eventsRoles = append(eventsRoles, models.EventRole{
			EventId:     event.ID.Hex(),
			Role:        role,
			Permissions: permissions,
		})
	}

	return eventsRoles, nil
}

func getVolunteerFromEvent(volunteers []models.Volunteer, userId string) (models.Volunteer, bool) {
	for _, volunteer := range volunteers {
		if volunteer.UserId == userId {
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
	event.Volunteers = make([]models.Volunteer, 0)
	event.Lockers = make([]models.Locker, 0)
	event.Products = make([]models.Product, 0)
	event.EventTickets = make([]models.EventTicket, 0)

	result, err := collection.InsertOne(ctx, event)
	id := ""
	if err == nil {
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			id = oid.Hex()
		}
	}
	return id, err
}

func AddEventTicketToEvent(eventId string, eventTicket models.EventTicket) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	code := uuid.New().String()
	eventTicket.EventTicketCode = code

	update := bson.M{"$push": bson.M{"eventTickets": eventTicket}}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	return err
}

func AddVolunteerToEvent(eventId string, volunteer models.Volunteer) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	update := bson.M{"$push": bson.M{"volunteers": volunteer}}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	return err
}

func AddLockerToEvent(eventId string, locker models.Locker) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	update := bson.M{"$push": bson.M{"lockers": locker}}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	return err
}

func AddRangeOfLockersToEvent(eventId string, start, end int, prefix, suffix string) ([]models.Locker, error) {
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
		lockers = append(lockers, models.Locker{LockerCode: lockerCode, UserId: ""})
	}

	// Préparer la mise à jour pour ajouter les casiers
	update := bson.M{
		"$push": bson.M{"lockers": bson.M{"$each": lockers}},
	}

	objId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		return nil, err
	}

	// Exécuter la mise à jour
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return nil, err
	}

	return lockers, nil
}

func AddProductToEvent(eventId string, product models.Product) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	code := uuid.New().String()
	product.ProductCode = code

	update := bson.M{"$push": bson.M{"products": product}}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	return err
}

// <=== Modif dans la BDD ===>

func UpdateEvent(eventId string, updateData bson.M) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updateData})
	return err
}

func UpdateLocker(eventId string, locker models.Locker) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	update := bson.M{"$set": bson.M{"lockers.$.user": locker.UserId}}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId, "lockers.lockerCode": locker.LockerCode}, update)
	return err
}

func UpdateProduct(eventId, productCode string, updateData bson.M) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	// Créer un nouveau bson.M pour la mise à jour avec l'opérateur $.
	setUpdateData := bson.M{}
	for key, value := range updateData {
		setUpdateData["products.$."+key] = value
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId, "products.code": productCode}, bson.M{"$set": setUpdateData})
	return err
}

func UpdateVolunteer(eventId string, volunteer models.Volunteer) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	_, err := collection.UpdateOne(ctx, bson.M{
		"_id":                objId,
		"volunteers.user_id": volunteer.UserId},
		bson.M{"$set": bson.M{
			"volunteers.$.permissions": volunteer.Permissions,
			"volunteers.$.isAdmin":     volunteer.IsAdmin,
		},
		},
	)
	return err
}

func UpdateEventTicket(eventId, eventTicketCode string, updateData bson.M) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	// Créer un nouveau bson.M pour la mise à jour avec l'opérateur $.
	setUpdateData := bson.M{}
	for key, value := range updateData {
		setUpdateData["tickets.$."+key] = value
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId, "eventTickets.EventTicketCode": eventTicketCode}, bson.M{"$set": setUpdateData})
	return err
}

// <=== Suppression dans la BDD ===>

func DeleteEvent(eventId string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": objId})
	return err
}

func DeleteEventTicket(eventId, eventTicketCode string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$pull": bson.M{"eventTickets": bson.M{"eventTicketCode": eventTicketCode}}})
	return err
}

func DeleteLocker(eventId, lockerCode string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$pull": bson.M{"lockers": bson.M{"lockerCode": lockerCode}}})
	return err
}

func DeleteAllLockers(eventId string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convertir l'ID de l'événement en ObjectID
	objId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		return err
	}

	// Mise à jour pour supprimer tous les casiers
	update := bson.M{
		"$set": bson.M{"lockers": []models.Locker{}}, // Définir les casiers à une liste vide
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return err
	}

	return nil
}

func DeleteVolunteer(eventId, userId string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		return err
	}

	// Vérifier si l'utilisateur est un administrateur de l'événement
	var event models.Event
	if err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&event); err != nil {
		return err
	}

	for _, volunteer := range event.Volunteers {
		if volunteer.UserId == userId && volunteer.IsAdmin {
			return fmt.Errorf("impossible de supprimer un bénévole qui est également administrateur")
		}
	}

	// Supprimer le bénévole si ce n'est pas un administrateur
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$pull": bson.M{"volunteers": bson.M{"userId": userId}}})
	return err
}

func DeleteProduct(eventId, productCode string) error {
	collection := Database.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$pull": bson.M{"products": bson.M{"code": productCode}}})
	return err
}
