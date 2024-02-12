package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	EventId         primitive.ObjectID `bson:"EventId"`
	UserId          primitive.ObjectID `bson:"UserId"`
	EventTicketCode string             `bson:"EventTicketCode"`
	IsSam           bool               `bson:"IsSam,omitempty"`
	IsUsed          bool               `bson:"IsUsed,omitempty"`
	WasPurchased    bool               `bson:"WasPurchased"`
	ReturnedEcoCup  bool               `bson:"ReturnedEcoCup,omitempty"`

	//Pour les informations simple sur un utilisateur
	User *User `bson:"User,omitempty"`
}

type TicketRequestData struct {
	UserId       string `json:"UserId"`
	WasPurchased *bool  `json:"WasPurchased,omitempty"`
	IsSam        bool   `json:"IsSam,omitempty"`
	IsUsed       bool   `json:"IsUsed,omitempty"`
}
