package database

import (
	"fmt"
	"log"
	"time"

	"github.com/bonsus/go-saas/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg config.DBConfig) (DB *gorm.DB, err error) {
	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.TZ,
	)
	DB, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return nil, err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}

	sqlDB.SetMaxIdleConns(int(cfg.ConnectionPool.MaxIdleConnection))                                // Misalnya, 10 koneksi idle maksimal
	sqlDB.SetMaxOpenConns(int(cfg.ConnectionPool.MaxOpenConnection))                                // Misalnya, 100 koneksi maksimal yang bisa dibuka
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnectionPool.MaxLifetimeConnection) * time.Second) // Tidak ada batas waktu koneksi
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnectionPool.MaxIdletimeConnection) * time.Second)

	fmt.Println("Connected to the database successfully!")
	return DB, nil
}
