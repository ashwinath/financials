package context

import (
	"fmt"

	"github.com/ashwinath/financials/api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Context struct {
	DB *gorm.DB
}

func InitContext(c *config.Config) (*Context, error) {
	context := Context{}
	db, err := initDB(c.Database)
	if err != nil {
		return nil, err
	}
	context.DB = db

	return &context, nil
}

func initDB(dbConfig config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port,
		dbConfig.TimeZone,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
