package postgres

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(AccessTokenEntity{})
	if err != nil {
		log.Printf("Error during run migration %v", err)
	}
}
