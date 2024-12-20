package migrate

import (
	"email-service/internal/models"
	"log"

	"gorm.io/gorm"
)

func ModelsAutoMigrate(db *gorm.DB) error {

	modelsToMigrate := []interface{}{
		&models.EmailTask{},
		&models.Attachment{},
	}

	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	log.Println("Database schema migrated successfully!")
	return nil
}
