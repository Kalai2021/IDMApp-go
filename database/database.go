package database

import (
	"fmt"
	"log"

	"idmapp-go/config"
	"idmapp-go/internal/client"
	"idmapp-go/internal/group"
	"idmapp-go/internal/member"
	"idmapp-go/internal/org"
	"idmapp-go/internal/pkce"
	"idmapp-go/internal/role"
	"idmapp-go/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	dsn := cfg.GetDatabaseDSN()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&user.User{},
		&group.Group{},
		&role.Role{},
		&org.Org{},
		&member.Member{},
		&pkce.PKCECode{},
		&client.Client{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
