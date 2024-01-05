package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Public      bool               `bson:"public,omitempty"`
	Description string             `bson:"description,omitempty"`
	Location    string             `bson:"location,omitempty"`
	Contact     string             `bson:"contact,omitempty"`
	ImageURL    string             `bson:"image_url,omitempty"`
	Tickets     []EventTicket      `bson:"tickets,omitempty"`
	Volunteers  []Volunteer        `bson:"volunteers,omitempty"`
	Admins      []string           `bson:"admins,omitempty"`
	Products    []Product          `bson:"products,omitempty"`
	Lockers     []Locker           `bson:"lockers,omitempty"`
}

type EventTicket struct {
	Code  string  `bson:"code,omitempty"`
	Title string  `bson:"title,omitempty"`
	Price float64 `bson:"price,omitempty"`
}

type Volunteer struct {
	User_ID     string   `bson:"user,omitempty"`        // ID de l'utilisateur
	Permissions []string `bson:"permissions,omitempty"` // Liste des permissions
}
type Product struct {
	Code        string  `bson:"code,omitempty"`
	Label       string  `bson:"label,omitempty"`
	Price       float64 `bson:"price,omitempty"`
	Is_Adult    bool    `bson:"adult,omitempty"` // Indique si le produit est pour adultes
	Icon        string  `bson:"icon,omitempty"`
	End_Of_Sale string  `bson:"end_of_sale,omitempty"` // Date/Heure de fin de vente
}
type Locker struct {
	Code string `bson:"code,omitempty"`
	User string `bson:"user"`
}

type EventRole struct {
	EventID string   `json:"event_id"`
	Role    string   `json:"role"`
	Rights  []string `json:"rights"`
}
