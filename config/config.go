package config

import (
	"fmt"
	"log"
	"mediadex/database"
	"mediadex/database/dbmodel"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Constants struct {
	Port             string `yaml:"port"`
	JWTSecret        string `yaml:"jwt_secret"`
	ConnectionString string `yaml:"connection_string"`
}

type Config struct {
	Constants

	// Repositories
	CollectionRepository dbmodel.CollectionRepository
	FieldRepository      dbmodel.FieldRepository
	MediaRepository      dbmodel.MediaRepository
	MediaFieldRepository dbmodel.MediaFieldRepository
	TagRepository        dbmodel.TagRepository
	UserRepository       dbmodel.UserRepository
}

func initEnv(fileName string) (Constants, error) {
	// Load .env file
	if err := godotenv.Load(fileName); err != nil {
		return Constants{}, fmt.Errorf("error loading env file %q: %w", fileName, err)
	}

	var constants Constants

	// Load constants from file
	constants.Port = os.Getenv("PORT")
	constants.JWTSecret = os.Getenv("JWT_SECRET_KEY")
	constants.ConnectionString = os.Getenv("CONNECTION_STRING")

	// Simple checks
	if constants.Port == "" {
		return Constants{}, fmt.Errorf("missing required env var PORT")
	}
	if constants.JWTSecret == "" {
		return Constants{}, fmt.Errorf("missing required env var JWT_SECRET_KEY")
	}
	if constants.ConnectionString == "" {
		return Constants{}, fmt.Errorf("missing required env var CONNECTION_STRING")
	}

	return constants, nil
}

func New() (*Config, error) {
	config := Config{}

	// Constants
	constants, err := initEnv(".env")

	config.Constants = constants
	if err != nil {
		return &config, err
	}

	// Open PostgreSQL Database Connection
	databaseSession, err := gorm.Open(postgres.Open(config.Constants.ConnectionString), &gorm.Config{})
	if err != nil {
		return &config, err
	}
	log.Println("Successfully connected to postgres database")

	err = database.Migrate(databaseSession)
	if err != nil {
		return &config, err
	}
	log.Println("Successfully migrated database")

	// Repositories Initialization
	config.CollectionRepository = dbmodel.NewCollectionRepository(databaseSession)
	config.FieldRepository = dbmodel.NewFieldRepository(databaseSession)
	config.MediaRepository = dbmodel.NewMediaRepository(databaseSession)
	config.MediaFieldRepository = dbmodel.NewMediaFieldRepository(databaseSession)
	config.TagRepository = dbmodel.NewTagRepository(databaseSession)
	config.UserRepository = dbmodel.NewUserRepository(databaseSession)

	return &config, nil
}
