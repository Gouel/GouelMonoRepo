package database

import (
	"context"
	"log"
	"os/exec"
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

func ImportDB(filePath string, cfg config.Config) error {

	cmd := exec.Command("mongorestore", "--uri="+cfg.MongoDBURI, "--gzip", "--archive="+filePath)
	if err := cmd.Run(); err != nil {
		log.Printf("Erreur lors de l'importation de la base de données: %v", err)
		return err
	}
	log.Println("Importation réussie")
	return nil
}
func ExportDB(filePath string, cfg config.Config) error {

	cmd := exec.Command("mongodump", "--uri="+cfg.MongoDBURI, "--db", cfg.MongoDBName, "--gzip", "--archive="+filePath)
	if err := cmd.Run(); err != nil {
		log.Printf("Erreur lors de l'exportation de la base de données: %v", err)
		return err
	}
	log.Printf("Exportation réussie vers %s", filePath)
	return nil
}
