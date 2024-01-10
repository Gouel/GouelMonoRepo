package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTicket(userId, eventId, eventTicketCode string, wasPurchased bool) (string, error) {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newTicket := models.Ticket{
		EventId:      eventId,
		IsSam:        false,
		IsUsed:       false,
		WasPurchased: wasPurchased,
	}

	result, err := collection.InsertOne(ctx, newTicket)
	if err != nil {
		return "", err
	}

	// Obtenez l'ID du nouveau ticket
	newTicketId := result.InsertedID.(primitive.ObjectID)
	return newTicketId.Hex(), nil
}

func ValidateTicket(ticketId, eventId string) (int, error) {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(ticketId)

	// Trouver le ticket et vérifier s'il est déjà validé
	var ticket models.Ticket
	err := collection.FindOne(ctx, bson.M{"_id": objId, "eventId": eventId}).Decode(&ticket)
	if err != nil {
		return 1, err // Retourne l'erreur si le ticket n'est pas trouvé
	}

	// Vérifier si le ticket est déjà validé
	if ticket.IsUsed {
		return 2, fmt.Errorf("le ticket avec l'ID %s est déjà validé", ticketId)
	}

	// Mettre à jour le ticket pour le valider
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"isUsed": true}})
	if result.Err() != nil {
		return 3, result.Err() // Retourne une erreur en cas de problème lors de la mise à jour
	}

	return 0, nil
}

func SetSAM(ticketId, eventId string, isSam bool) error {
	collection := Database.Collection("tickets")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(ticketId)
	// Mettre à jour le ticket pour le valider
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": objId, "eventId": eventId}, bson.M{"$set": bson.M{"isSam": isSam}})
	if result.Err() != nil {
		return result.Err() // Retourne une erreur en cas de problème lors de la mise à jour
	}
	return nil
}

func DeleteTicket(ticketId string) error {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(ticketId)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": objId})
	return err
}

func GetTicketInfo(ticketId string, eventId *string) (*models.Ticket, error) {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(ticketId)

	var ticket models.Ticket

	var filter bson.M

	if eventId != nil {
		filter = bson.M{"_id": objId, "eventId": eventId}
	} else {
		filter = bson.M{"_id": objId}
	}

	err := collection.FindOne(ctx, filter).Decode(&ticket)
	if err != nil {
		return nil, err
	}

	user, err := userFromTicket(ctx, ticket)
	if err != nil {
		return nil, err
	}
	ticket.User = user
	return &ticket, nil
}

func userFromTicket(ctx context.Context, ticket models.Ticket) (*models.User, error) {
	var user models.User
	userId, _ := primitive.ObjectIDFromHex(ticket.UserId)
	err := Database.Collection("users").FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllTicketsFromEvent(eventId string) (bson.A, error) {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	tickets := make(bson.A, 0)

	for cursor.Next(ctx) {
		var ticket models.Ticket
		if err := cursor.Decode(&ticket); err != nil {
			continue
		}
		user, err := userFromTicket(ctx, ticket)
		if err != nil {
			continue
		}
		ticket.User = user
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}
