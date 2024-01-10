package database

import (
	"context"
	"log"
	"time"

	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func ConnectDB(cfg config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBURI))
	if err != nil {
		log.Fatal(err)
	}

	// Vérifier la connexion
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connexion à MongoDB réussie.")
	Client = client
	Database = client.Database(cfg.MongoDBName)
}
