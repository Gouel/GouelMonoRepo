package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config représente la structure de la configuration de votre serveur
type Config struct {
	MongoDBURI     string
	MongoDBName    string
	JWTSecretKey   string
	JWTExpiration  int
	ServerPort     string
	ServerHost     string
	DebugMode      bool
	TrustedProxies []string

	//config SMTP
	SMTPServer    string
	SMTPPort      string
	Email         string
	EmailPassword string
	SMTPUseSSL    bool
}

// LoadConfig charge la configuration à partir des variables d'environnement
func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		serverPort = "8080" // Port par défaut
	}

	serverHost, exists := os.LookupEnv("SERVER_HOST")
	if !exists {
		serverHost = "127.0.0.1" // localhost par défaut
	}

	debugMode, err := strconv.ParseBool(os.Getenv("DEBUG_MODE"))
	if err != nil {
		debugMode = false // Valeur par défaut si non définie ou invalide
	}

	trustedProxies := strings.Split(os.Getenv("TRUSTED_PROXIES"), ",")

	expiration, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))

	if err != nil {
		log.Fatalf("Error JWTExpiration : %v", err)
	}

	// Config SMTP

	useSSL, err := strconv.ParseBool(os.Getenv("USE_SSL"))
	if err != nil {
		debugMode = false // Valeur par défaut si non définie ou invalide
	}
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("EMAIL")
	smtpPassword := os.Getenv("EMAIL_PASSWORD")

	return Config{
		MongoDBURI:     os.Getenv("MONGODB_URI"),
		MongoDBName:    os.Getenv("MONGODB_DB_NAME"),
		JWTSecretKey:   os.Getenv("JWT_SECRET_KEY"),
		JWTExpiration:  expiration,
		ServerPort:     serverPort,
		ServerHost:     serverHost,
		DebugMode:      debugMode,
		TrustedProxies: trustedProxies,

		SMTPServer:    smtpServer,
		SMTPPort:      smtpPort,
		SMTPUseSSL:    useSSL,
		Email:         smtpEmail,
		EmailPassword: smtpPassword,
	}, nil
}
