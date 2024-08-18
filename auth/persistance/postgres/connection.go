package postgres

import (
	log "github.com/sirupsen/logrus"
	psq "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"my-stocks/auth/config"
)

var db *gorm.DB

func GetConnection(cfg config.Database) *gorm.DB {
	if db != nil {
		return db
	}
	conn, err := gorm.Open(psq.New(psq.Config{
		DSN:                  "host=" + cfg.Host + " user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.Name + " port=" + cfg.Port + " sslmode=disable TimeZone=UTC",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db = conn

	return db
}
