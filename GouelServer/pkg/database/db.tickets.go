package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateTicket(userId, eventId, eventTicketCode string, ticketRequestData models.TicketRequestData) (string, error) {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIdOID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", err
	}
	eventIdOID, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		return "", err
	}
	newTicket := models.Ticket{
		EventId:         eventIdOID,
		EventTicketCode: eventTicketCode,
		IsSam:           ticketRequestData.IsSam,
		IsUsed:          ticketRequestData.IsUsed,
		WasPurchased:    *ticketRequestData.WasPurchased,
		UserId:          userIdOID,
	}

	result, err := collection.InsertOne(ctx, newTicket)
	if err != nil {
		return "", err
	}

	//Ajout du SamBonus si le ticket est un Sam
	if ticketRequestData.IsSam {
		user, err := GetUserById(userIdOID)
		if err != nil {
			return "", err
		}
		event, err := GetEventById(eventIdOID.Hex())
		if err != nil {
			return "", err
		}
		solde := user.Solde + event.SamBonus
		err = UpdateUser(userId, bson.M{"Solde": solde})
		if err != nil {
			return "", err
		}
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
	err := collection.FindOne(ctx, bson.M{"_id": objId, "EventId": eventId}).Decode(&ticket)
	if err != nil {
		return 1, err // Retourne l'erreur si le ticket n'est pas trouvé
	}

	// Vérifier si le ticket est déjà validé
	if ticket.IsUsed {
		return 2, fmt.Errorf("le ticket avec l'ID %s est déjà validé", ticketId)
	}

	// Mettre à jour le ticket pour le valider
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"IsUsed": true}})
	if result.Err() != nil {
		return 3, result.Err() // Retourne une erreur en cas de problème lors de la mise à jour
	}

	return 0, nil
}

func SetSAM(ticketId, eventId string, isSam bool) error {
	collection := Database.Collection("tickets")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println(ticketId, eventId)
	ticket, err := GetTicketInfo(ticketId, &eventId)
	if err != nil {
		fmt.Println("A")
		return err
	}
	// Mettre à jour le ticket pour le valider
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": ticket.ID, "EventId": ticket.EventId}, bson.M{"$set": bson.M{"IsSam": isSam}})
	if result.Err() != nil {
		fmt.Println("Z")
		return result.Err() // Retourne une erreur en cas de problème lors de la mise à jour
	}

	user, err := GetUserById(ticket.UserId)
	if err != nil {
		fmt.Println("B")
		return err
	}
	event, err := GetEventById(ticket.EventId.Hex())
	if err != nil {
		fmt.Println("C")
		return err
	}

	var solde float32

	if isSam {
		//On ajoute le SamBonus à l'utilisateur
		solde = user.Solde + event.SamBonus
	} else {
		if user.Solde > event.SamBonus {
			solde = user.Solde - event.SamBonus
		}
	}

	err = UpdateUser(ticket.UserId.Hex(), bson.M{"Solde": solde})
	if err != nil {
		fmt.Println("D")
		return err
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
		eventIdOID, err := primitive.ObjectIDFromHex(*eventId)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"_id": objId, "EventId": eventIdOID}
	} else {
		filter = bson.M{"_id": objId}
	}

	err := collection.FindOne(ctx, filter).Decode(&ticket)
	if err != nil {
		fmt.Printf("1 debug : %s", err)
		return nil, err
	}

	user, err := userFromTicket(ctx, ticket)
	if err != nil {
		fmt.Printf("2 debug : %s", err)
		return nil, err
	}
	ticket.User = user
	return &ticket, nil
}

func userFromTicket(ctx context.Context, ticket models.Ticket) (*models.User, error) {
	var user models.User
	err := Database.Collection("users").FindOne(
		ctx, bson.M{"_id": ticket.UserId},
		options.FindOne().SetProjection(bson.M{
			"FirstName": 1,
			"LastName":  1,
			"DOB":       1,
			"Email":     1,
		}),
	).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllTicketsFromEvent(eventId string) (bson.A, error) {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventIdOID, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		return make(primitive.A, 0), err
	}

	filter := bson.M{"EventId": eventIdOID}
	cursor, err := collection.Find(ctx, filter)
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

func GetPaginatedTicketsFromEvent(eventId string, page int64) (bson.M, error) {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//On skip les tickets des pages précédentes
	opts := options.Find().SetSkip(25 * (page - 1)).SetLimit(25) // 25 tickets par page

	filter := bson.M{"EventId": eventId}

	cursor, err := collection.Find(ctx, filter, opts)
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

	opts2 := options.Count()
	count, err := collection.CountDocuments(ctx, filter, opts2)
	if err != nil {
		return nil, err
	}

	return bson.M{
		"Tickets": tickets,
		"Total":   count,
		"PerPage": 25,
	}, nil
}

func ReturnEcoCup(ticketId string) error {
	collection := Database.Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(ticketId)

	// Trouver le ticket et vérifier s'il est déjà validé
	var ticket models.Ticket
	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&ticket)
	if err != nil {
		return err // Retourne l'erreur si le ticket n'est pas trouvé
	}

	// Vérifier une écocup a déjà été rendue
	if ticket.ReturnedEcoCup {
		return errors.New("an eco-cup has already been returned for this ticket")
	}

	// Mettre à jour le ticket pour le valider
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"ReturnedEcoCup": true}})
	if result.Err() != nil {
		return result.Err() // Retourne une erreur en cas de problème lors de la mise à jour
	}

	return nil
}
