package database

import (
	"MerchShop/entities"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var predefinedMerchItems = []entities.Merch{
	{Type: "t-shirt", Cost: 80},
	{Type: "cup", Cost: 20},
	{Type: "book", Cost: 50},
	{Type: "pen", Cost: 10},
	{Type: "powerbank", Cost: 200},
	{Type: "hoody", Cost: 300},
	{Type: "umbrella", Cost: 200},
	{Type: "socks", Cost: 10},
	{Type: "wallet", Cost: 50},
	{Type: "pink-hoody", Cost: 500},
}

var Instance *gorm.DB
var err error

func Connect(connectionString string) {
	Instance, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database...")
}

func Migrate() {
	Instance.AutoMigrate(&entities.User{})
	Instance.AutoMigrate(&entities.Merch{})
	Instance.AutoMigrate(&entities.Transaction{})
	Instance.AutoMigrate(&entities.Purchase{})
	log.Println("Database Migration Completed...")

	InitializeMerchTable(Instance)
	log.Println("Merch table is ready")
}

func InitializeMerchTable(db *gorm.DB) {
	for _, item := range predefinedMerchItems {
		var existing entities.Merch
		if err := db.Where("type = ?", item.Type).First(&existing).Error; err != nil {
			db.Create(&item)
		}
	}
}
