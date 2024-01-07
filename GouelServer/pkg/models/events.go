package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Title        string             `bson:"title,omitempty"`
	IsPublic     bool               `bson:"isPublic,omitempty"`
	Description  string             `bson:"description,omitempty"`
	Location     string             `bson:"location,omitempty"`
	Contact      string             `bson:"contact,omitempty"`
	EventTickets []EventTicket      `bson:"eventTickets,omitempty"`
	Volunteers   []Volunteer        `bson:"volunteers,omitempty"`
	Products     []Product          `bson:"products,omitempty"`
	Lockers      []Locker           `bson:"lockers,omitempty"`
}

type EventTicket struct {
	EventTicketCode string  `bson:"EventTicketCode,omitempty"`
	Title           string  `bson:"title,omitempty"`
	Price           float64 `bson:"price,omitempty"`
}

type Volunteer struct {
	UserId      string   `bson:"userId,omitempty"`      // ID de l'utilisateur
	Permissions []string `bson:"permissions,omitempty"` // Liste des permissions
	IsAdmin     bool     `bson:"isAdmin, omitempty"`
}
type Product struct {
	ProductCode string  `bson:"productCode,omitempty"`
	Label       string  `bson:"label,omitempty"`
	Price       float32 `bson:"price,omitempty"`
	HasAlcohol  bool    `bson:"hasAlcohol,omitempty"` // Indique si le produit est pour adultes
	Icon        string  `bson:"icon,omitempty"`
	EndOfSale   string  `bson:"endOfSale,omitempty"` // Date/Heure de fin de vente
	Amount      *int    `bson:"amount,omitempty"`
	Purchased   int     `bson:"purchased,omitempty"`
}

type Locker struct {
	LockerCode string `bson:"lockerCode,omitempty"`
	UserId     string `bson:"userId"`
}

type EventRole struct {
	EventId     string   `bson:"eventId"`
	Role        string   `bson:"role"`
	Permissions []string `bson:"permissions,omitempty"`
}
