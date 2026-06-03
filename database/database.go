package database

import (
	"fmt"
	"mediadex/database/dbmodel"

	"gorm.io/gorm"
)

func Migrate(database *gorm.DB) error {
	err := database.AutoMigrate(
		&dbmodel.Collection{},
		&dbmodel.Field{},
		&dbmodel.Media{},
		&dbmodel.MediaField{},
		&dbmodel.Tag{},
		&dbmodel.User{},
	)

	if err != nil {
		return fmt.Errorf("Failed to migrate database: %w", err)
	}

	return nil
}
