package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	EventId         string             `bson:"eventId"`
	UserId          string             `bson:"userId"`
	EventTicketCode string             `bson:"eventTicketCode"`
	IsSam           bool               `bson:"isSam,omitempty"`
	IsUsed          bool               `bson:"isUsed,omitempty"`
	WasPurchased    bool               `bson:"wasPurchased"`

	//Pour les informations simple sur un utilisateur
	User *User `bson:"user,omitempty"`
}
