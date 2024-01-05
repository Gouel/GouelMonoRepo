package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTicket(userID, eventID, ticketCode string, isSpecial bool) (string, error) {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newTicket := bson.M{
		"user": userID,
		"event": bson.M{
			"event_id":    eventID,
			"ticket_code": ticketCode,
		},
		"valid":   false,
		"SAM":     false,
		"special": isSpecial,
	}

	result, err := collection.InsertOne(ctx, newTicket)
	if err != nil {
		return "", err
	}

	// Obtenez l'ID du nouveau ticket
	newTicketID := result.InsertedID.(primitive.ObjectID)
	return newTicketID.Hex(), nil
}

func ValidateTicket(ticketID, eventID string) (int, error) {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(ticketID)

	// Trouver le ticket et vérifier s'il est déjà validé
	var ticket bson.M
	err := collection.FindOne(ctx, bson.M{"_id": objID, "event.event_id": eventID}).Decode(&ticket)
	if err != nil {
		return 1, err // Retourne l'erreur si le ticket n'est pas trouvé
	}

	// Vérifier si le ticket est déjà validé
	if ticket["valid"].(bool) {
		return 2, fmt.Errorf("le ticket avec l'ID %s est déjà validé", ticketID)
	}

	// Mettre à jour le ticket pour le valider
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"valid": true}})
	if result.Err() != nil {
		return 3, result.Err() // Retourne une erreur en cas de problème lors de la mise à jour
	}

	return 0, nil
}

func SetSAM(ticketID, eventID string, isSAM bool) error {
	collection := Database.Collection("tickets")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objID, _ := primitive.ObjectIDFromHex(ticketID)
	// Mettre à jour le ticket pour le valider
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": objID, "event.event_id": eventID}, bson.M{"$set": bson.M{"SAM": isSAM}})
	if result.Err() != nil {
		return result.Err() // Retourne une erreur en cas de problème lors de la mise à jour
	}
	return nil
}

func DeleteTicket(ticketID string) error {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(ticketID)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func GetTicketInfo(ticketID string, eventID *string) (bson.M, error) {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(ticketID)

	var ticket bson.M

	var filter bson.M

	if eventID != nil {
		filter = bson.M{"_id": objID, "event.event_id": eventID}
	} else {
		filter = bson.M{"_id": objID}
	}

	err := collection.FindOne(ctx, filter).Decode(&ticket)
	if err != nil {
		return nil, err
	}

	user, err := userFromTicker(ctx, ticket)
	if err != nil {
		return nil, err
	}
	ticket["user"] = user
	return ticket, nil
}

func userFromTicker(ctx context.Context, ticket primitive.M) (bson.M, error) {
	var user bson.M
	userId, _ := primitive.ObjectIDFromHex(ticket["user"].(string))
	err := Database.Collection("users").FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	userReturn := bson.M{
		"firstName": user["firstName"],
		"lastName":  user["lastName"],
		"email":     user["email"],
		"dob":       user["dob"],
		"user_id":   user["_id"].(primitive.ObjectID).Hex(),
	}

	return userReturn, nil
}

func GetAllTicketsFromEvent(eventID string) (bson.A, error) {
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
		var ticket bson.M
		if err := cursor.Decode(&ticket); err != nil {
			continue
		}
		user, err := userFromTicker(ctx, ticket)
		if err != nil {
			continue
		}
		ticket["user"] = user
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}
