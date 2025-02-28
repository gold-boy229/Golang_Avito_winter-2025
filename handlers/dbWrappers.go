package handlers

import (
	db "MerchShop/database"
	"MerchShop/entities"
	"MerchShop/model"

	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	return db.Instance
}

func getMerchByType(merchType string) (merch entities.Merch, err error) {
	return db.GetMerchByType(merchType)
}

func createNewUser(newUser entities.User) (entities.User, error) {
	return db.CreateNewUser(newUser)
}

func getUserByUsername(username string) (user entities.User, err error) {
	return db.GetUserByUsername(username)
}

func getUserInventoryItems(user entities.User) (inventoryItems []model.InventoryItem, err error) {
	return db.GetUserInventoryItems(user)
}

func getReceivedOperations(user entities.User) (receivedOperations []model.ReceivedOperation, err error) {
	return db.GetReceivedOperations(user)
}

func getSentOperations(user entities.User) (sentOperations []model.SentOperation, err error) {
	return db.GetSentOperations(user)
}
