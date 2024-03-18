package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/config"
	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/database"
	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/models"
	routes "github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Charger la configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la configuration: %v", err)
	}

	database.ConnectDB(cfg)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--setup":
			setupSuperAdmin()
			setupAPI()
		case "--secret":
			// 64 bits peu sécurisé => 256bits OK
			generateSecretKey(256)
		case "--import":
			if len(os.Args) != 3 {
				log.Fatalf("Usage: gouel --import /path/to/the/in.dump")
			}
			importPath := os.Args[2]
			err := database.ImportDB(importPath, cfg)
			if err != nil {
				log.Fatalf("Erreur lors de l'importation de la base de données: %v", err)
			}
		case "--export":
			if len(os.Args) != 3 {
				log.Fatalf("Usage: gouel --export /path/to/the/out.dump")
			}
			exportPath := os.Args[2]
			err := database.ExportDB(exportPath, cfg)
			if err != nil {
				log.Fatalf("Erreur lors de l'exportation de la base de données: %v", err)
			}

		case "--help":
			//Affichage aide / usage
			fmt.Println("<=== Gouel ===>")
			aide := []string{
				"--export <file out>",
				"--import <file in>",
				"--help",
				"--secret",
				"--setup",
			}

			for _, v := range aide {
				fmt.Println("gouel " + v)
			}

			return

		default:
			//Affichage erreur usage
			log.Println("Argument non reconnu !")
			log.Println("Usage: gouel --help")
			return
		}
		return
	}

	// Configurer le mode de débogage de Gin
	if cfg.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialiser le framework Gin
	router := gin.Default()

	// Configurer Gin pour envoyer les rêquetes CORS
	router.Use(cors.Default())

	// Configurer Gin pour ne pas faire confiance à tous les proxies
	router.SetTrustedProxies(cfg.TrustedProxies)

	routes.Routes(router, cfg)

	// Lancement du serveur sur le port configuré
	// err = router.RunTLS(cfg.ServerHost+":"+cfg.ServerPort, "server.pem", "key.pem")
	err = router.Run(cfg.ServerHost + ":" + cfg.ServerPort)
	if err != nil {
		log.Fatalf("Erreur lors du lancement du serveur: %v", err)
	}

	fmt.Printf("Serveur lancé sur l'adresse %s\n", cfg.ServerHost+":"+cfg.ServerPort)
}

func setupSuperAdmin() {
	password := os.Getenv("SUPERADMIN_PASSWORD")
	email := os.Getenv("SUPERADMIN_EMAIL")
	dob := os.Getenv("SUPERADMIN_DOB")

	superAdmin := models.User{
		FirstName: "SUPERADMIN",
		LastName:  "SUPERADMIN",
		Email:     email,
		DOB:       dob,
		Password:  password,
		Role:      "SUPERADMIN",
		Solde:     0,
	}

	_, err := database.CreateUser(superAdmin)
	if err != nil {
		log.Fatalf("Erreur lors de la création du super administrateur: %v", err)
	}

	log.Println("Super administrateur créé avec succès")
}

func setupAPI() {
	password := os.Getenv("API_PASSWORD")
	email := os.Getenv("API_EMAIL")
	userAPI := models.User{
		FirstName: "API",
		LastName:  "API",
		Email:     email,
		DOB:       "2000-01-01",
		Password:  password,
		Role:      "API",
		Solde:     0,
	}

	_, err := database.CreateUser(userAPI)
	if err != nil {
		log.Fatalf("Erreur lors de la création de l'utilisateur API: %v", err)
	}

	log.Println("Utilisateur API créé avec succès")
}

func generateSecretKey(length int) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	secretKey := base64.URLEncoding.EncodeToString(bytes)
	fmt.Printf(`JWT_SECRET_KEY="%s"`, secretKey)
}
