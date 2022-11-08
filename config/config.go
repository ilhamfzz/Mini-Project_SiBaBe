package config

import (
	"Mini-Project_SiBaBe/model"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost,
		dbUser,
		dbPass,
		dbName,
		dbPort,
	)
	var err error
	model.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return model.DB
}

func InitDatabaseTest() *gorm.DB {
	dbHostTest := os.Getenv("DB_HOST_TEST")
	dbUserTest := os.Getenv("DB_USER_TEST")
	dbPassTest := os.Getenv("DB_PASS_TEST")
	dbNameTest := os.Getenv("DB_NAME_TEST")
	dbPortTest := os.Getenv("DB_PORT_TEST")
	dsnTest := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHostTest,
		dbUserTest,
		dbPassTest,
		dbNameTest,
		dbPortTest,
	)
	var err error
	model.DB, err = gorm.Open(postgres.Open(dsnTest), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	model.DB.Migrator().DropTable(
		&model.Customer{}, &model.Admin{}, &model.Produk{},
		&model.Keranjang{}, &model.Produk_Keranjang{}, &model.Produksi{},
		&model.Produk_Produksi{}, &model.Pemesanan{}, &model.Admin_Pemesanan{},
		&model.Feedback_Pemesanan{}, &model.Feedback{},
	)
	model.DB.AutoMigrate(
		&model.Customer{}, &model.Admin{}, &model.Produk{},
		&model.Keranjang{}, &model.Produk_Keranjang{}, &model.Produksi{},
		&model.Produk_Produksi{}, &model.Pemesanan{}, &model.Admin_Pemesanan{},
		&model.Feedback_Pemesanan{}, &model.Feedback{},
	)

	return model.DB
}
