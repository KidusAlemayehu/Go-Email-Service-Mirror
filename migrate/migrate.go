package migrate

import (
	"email-service/internal/models"
	"email-service/utils/log"

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

	log.Logger.Info("Database schema migrated successfully!")
	return nil
}
