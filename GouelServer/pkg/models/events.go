package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Title        string             `bson:"Title,omitempty"`
	IsPublic     bool               `bson:"IsPublic,omitempty"`
	Description  string             `bson:"Description,omitempty"`
	Location     string             `bson:"Location,omitempty"`
	Contact      string             `bson:"Contact,omitempty"`
	EventTickets []EventTicket      `bson:"EventTickets,omitempty"`
	Volunteers   []Volunteer        `bson:"Volunteers,omitempty"`
	Products     []Product          `bson:"Products,omitempty"`
	Lockers      []Locker           `bson:"Lockers,omitempty"`
	Options      bson.M             `bson:"Options,omitempty"`
}

type EventTicket struct {
	EventTicketCode string  `bson:"EventTicketCode,omitempty"`
	Title           string  `bson:"Title,omitempty"`
	Price           float64 `bson:"Price,omitempty"`
	Amount          *int32  `bson:"Amount, omitempty"`
	Purchased       int32   `bson:"Purchased, omitempty"`
}

type Volunteer struct {
	UserId      string   `bson:"UserId,omitempty"`      // ID de l'utilisateur
	Permissions []string `bson:"Permissions,omitempty"` // Liste des permissions
	IsAdmin     bool     `bson:"IsAdmin, omitempty"`
}
type Product struct {
	ProductCode string  `bson:"ProductCode,omitempty"`
	Label       string  `bson:"Label,omitempty"`
	Price       float32 `bson:"Price,omitempty"`
	HasAlcohol  bool    `bson:"HasAlcohol,omitempty"` // Indique si le produit est pour adultes
	Icon        string  `bson:"Icon,omitempty"`
	EndOfSale   string  `bson:"EndOfSale,omitempty"` // Date/Heure de fin de vente
	Amount      *int    `bson:"Amount,omitempty"`
	Purchased   int     `bson:"Purchased,omitempty"`
}

type Locker struct {
	LockerCode string `bson:"LockerCode,omitempty"`
	UserId     string `bson:"UserId"`
}

type EventRole struct {
	EventId     string   `bson:"EventId"`
	Role        string   `bson:"Role"`
	Permissions []string `bson:"Permissions,omitempty"`
}
