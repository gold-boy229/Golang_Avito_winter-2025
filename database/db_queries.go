package database

import (
	"MerchShop/entities"
	"MerchShop/model"
	"fmt"
)

func GetMerchByType(merchType string) (merch entities.Merch, err error) {
	result := Instance.Where("type = ?", merchType).First(&merch)
	if result.Error != nil {
		err = fmt.Errorf("there is no merch with type %s", merchType)
	}
	return merch, err
}

func GetUserByUsername(username string) (user entities.User, err error) {
	result := Instance.Where("username = ?", username).First(&user)
	if result.Error != nil {
		err = fmt.Errorf("there is no user with username %s", username)
	}
	return
}

func GetUserInventoryItems(user entities.User) (inventoryItems []model.InventoryItem, err error) {
	result := Instance.Table("purchase").
		Select("merch.type, SUM(purchase.quantity) as quantity").
		Joins("INNER JOIN merch ON purchase.merchId = merch.Id").
		Where("purchase.user_id = ?", user.Id).
		Group("merch.type").
		Scan(&inventoryItems)

	if result.Error != nil {
		err = result.Error
	}
	return inventoryItems, err
}

func GetReceivedOperations(user entities.User) (receivedOperations []model.ReceivedOperation, err error) {
	result := Instance.Table("transaction").
		Select("from_user_id AS fromUserId, SUM(amount)").
		Where("to_user_id = ?", user.Id).
		Group("from_user_id").
		Scan(&receivedOperations)

	if result.Error != nil {
		err = result.Error
	}
	return receivedOperations, err
}

func GetSentOperations(user entities.User) (sentOperations []model.SentOperation, err error) {
	result := Instance.Table("transaction").
		Select("to_user_id AS toUserId, SUM(amount)").
		Where("from_user_id = ?", user.Id).
		Group("to_user_id").
		Scan(&sentOperations)

	if result.Error != nil {
		err = result.Error
	}
	return sentOperations, err
}
