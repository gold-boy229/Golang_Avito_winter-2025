package controllers

import (
	"MerchShop/database"
	"MerchShop/entities"
	"MerchShop/jwtutil"
	"MerchShop/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	curUser, err := getUserAfterMiddleware(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var (
		coins              uint = curUser.Balance
		inventoryItems     []model.InventoryItem
		receivedOperations []model.ReceivedOperation
		sentOperations     []model.SentOperation
	)

	err = getInventoryItems(&curUser, &inventoryItems)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = getReceivedOperations(&curUser, &receivedOperations)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = getSentOperations(&curUser, &sentOperations)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	coinHistory := model.CoinHistory{Received: receivedOperations, Sent: sentOperations}
	infoResponse := model.InfoResponse{Coins: coins, Inventory: inventoryItems, CoinHistory: coinHistory}

	respondJSON(w, http.StatusOK, infoResponse)
}

func getInventoryItems(user *entities.User, inventoryItems *[]model.InventoryItem) (err error) {
	err = database.Instance.Table("purchase").
		Select("merch.type, SUM(purchase.quantity) as quantity").
		Joins("INNER JOIN merch ON purchase.merchId = merch.Id").
		Where("purchase.user_id = ?", user.Id).
		Group("merch.type").
		Scan(&inventoryItems).Error
	return err
}

func getReceivedOperations(user *entities.User, receivedOperations *[]model.ReceivedOperation) (err error) {
	err = database.Instance.Table("transaction").
		Select("from_user_id AS fromUserId, SUM(amount)").
		Where("to_user_id = ?", user.Id).
		Group("from_user_id").
		Scan(&receivedOperations).Error
	return err
}

func getSentOperations(user *entities.User, sentOperations *[]model.SentOperation) (err error) {
	err = database.Instance.Table("transaction").
		Select("to_user_id AS toUserId, SUM(amount)").
		Where("from_user_id = ?", user.Id).
		Group("to_user_id").
		Scan(&sentOperations).Error
	return err
}

////

func SendCoinHandler(w http.ResponseWriter, r *http.Request) {
	var sendCoinRequest model.SendCoinRequest
	err := json.NewDecoder(r.Body).Decode(&sendCoinRequest)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body.\n"+err.Error())
		return
	}

	curUser, err := getUserAfterMiddleware(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	toUser, err := GetUserByUsername(sendCoinRequest.ToUser)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := sendCoinsToUser(&curUser, &toUser, sendCoinRequest.Amount)
	if err != nil {
		respondError(w, statusCode, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func sendCoinsToUser(fromUser *entities.User, toUser *entities.User, amount uint) (statusCode int, err error) {
	err = fromUser.CanSpendCoins(amount)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = fromUser.IsSameUser(toUser.Username)
	if err == nil {
		return http.StatusBadRequest, err
	}

	err = database.Instance.Transaction(func(tx *gorm.DB) error {
		fromUser.Balance -= amount
		if err := tx.Save(&fromUser).Error; err != nil {
			return err
		}

		toUser.Balance += amount
		if err := tx.Save(&toUser).Error; err != nil {
			return err
		}

		transaction := entities.Transaction{FromUserId: fromUser.Id, ToUserId: toUser.Id, Amount: amount}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

////

func BuyItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var merchType string = vars["item"]
	merch, err := GetMerchByType(merchType)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := getUserAfterMiddleware(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := buyMerch(&user, &merch)
	if err != nil {
		respondError(w, statusCode, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func buyMerch(user *entities.User, merch *entities.Merch) (statusCode int, err error) {
	err = user.CanSpendCoins(merch.Cost)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = database.Instance.Transaction(func(tx *gorm.DB) error {
		purchase := entities.Purchase{UserId: user.Id, MerchId: merch.Id}
		if err := tx.Create(&purchase).Error; err != nil {
			return err
		}

		user.Balance -= merch.Cost
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

////

func GetAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	var authRequest model.AuthRequest

	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body.\n"+err.Error())
		return
	}

	if authRequest.Username == "" || authRequest.Password == "" {
		respondError(w, http.StatusBadRequest, "username and password shouldn't be empty")
		return
	}

	// Check if user exists in BD:
	//   yes) check password and give token
	//   no) create newUser and give token
	user, err := GetUserByUsername(authRequest.Username)
	if err == nil {
		if err := user.CheckPassword(authRequest.Password); err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
	} else {
		newUser := entities.User{Username: authRequest.Username, Password: authRequest.Password}
		if err := database.Instance.Create(&newUser).Error; err != nil {
			respondError(w, http.StatusInternalServerError, "Couldn't create a new user.\n"+err.Error())
			return
		}
	}

	tokenJWT, err := jwtutil.GenerateJWT(authRequest.Username)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, model.AuthResponse{Token: tokenJWT})
}

////

func getUserAfterMiddleware(r *http.Request) (user entities.User, err error) {
	username := r.Header.Get(USERNAME)
	return GetUserByUsername(username)
}

func GetUserByUsername(username string) (user entities.User, err error) {
	user = entities.User{}
	err = nil
	result := database.Instance.Where("username = ?", username).First(&user)
	if result.Error != nil {
		err = fmt.Errorf("there is no user with username %s", username)
	}
	return
}

func GetMerchByType(merchType string) (merch entities.Merch, err error) {
	merch = entities.Merch{}
	err = nil
	result := database.Instance.Where("type = ?", merchType).First(&merch)
	if result.Error != nil {
		err = fmt.Errorf("there is no merch with type %s", merchType)
	}
	return
}

////

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusNotFound, "Wrong path or this function isn't implemented yet")
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}
