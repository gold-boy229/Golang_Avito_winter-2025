package database

import (
	"MerchShop/entities"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Instance *gorm.DB

func Connect(connectionString string) {
	var err error
	Instance, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database...")
}

func Migrate() {
	runDatabaseMigrations()
	log.Println("Database Migration Completed...")

	initializeMerchTable()
	log.Println("Merch table is ready")
}

func runDatabaseMigrations() {
	Instance.AutoMigrate(&entities.User{})
	Instance.AutoMigrate(&entities.Merch{})
	Instance.AutoMigrate(&entities.Transaction{})
	Instance.AutoMigrate(&entities.Purchase{})
}

func initializeMerchTable() {
	predefinedMerchItems := getPredefinedMerchItems()
	for _, merchItem := range predefinedMerchItems {
		if merchExists(merchItem) {
			insertMerchItemToDB(merchItem)
		}
	}
}

func getPredefinedMerchItems() []entities.Merch {
	return []entities.Merch{
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
}

func merchExists(merch entities.Merch) bool {
	_, err := GetMerchByType(merch.Type)
	return (err == nil)
}

func insertMerchItemToDB(merchItem entities.Merch) {
	Instance.Create(&merchItem)
}
