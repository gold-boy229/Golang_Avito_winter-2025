package controllers

import (
	"MerchShop/database"
	"MerchShop/entities"
	"MerchShop/jwtutil"
	"MerchShop/model"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	user, err := database.GetUserByUsername(getUsernameFromHeader(r))
	if err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	var (
		coins              uint = user.Balance
		inventoryItems     []model.InventoryItem
		receivedOperations []model.ReceivedOperation
		sentOperations     []model.SentOperation
	)

	inventoryItems, err = database.GetUserInventoryItems(user)
	if err != nil {
		respondInternalServerError(w, err.Error())
		return
	}

	receivedOperations, err = database.GetReceivedOperations(user)
	if err != nil {
		respondInternalServerError(w, err.Error())
		return
	}

	sentOperations, err = database.GetSentOperations(user)
	if err != nil {
		respondInternalServerError(w, err.Error())
		return
	}

	coinHistory := model.CoinHistory{Received: receivedOperations, Sent: sentOperations}
	infoResponse := model.InfoResponse{Coins: coins, Inventory: inventoryItems, CoinHistory: coinHistory}

	respondJSON(w, http.StatusOK, infoResponse)
}

////

func SendCoinHandler(w http.ResponseWriter, r *http.Request) {
	sendCoinRequest, err := getSendCoinRequest(r)
	if err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	curUser, err := database.GetUserByUsername(getUsernameFromHeader(r))
	if err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	toUser, err := database.GetUserByUsername(sendCoinRequest.ToUser)
	if err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	statusCode, err := sendCoinsToUser(curUser, toUser, sendCoinRequest.Amount)
	if err != nil {
		respondError(w, statusCode, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func getSendCoinRequest(r *http.Request) (sendCoinRequest model.SendCoinRequest, err error) {
	err = json.NewDecoder(r.Body).Decode(&sendCoinRequest)
	if err != nil {
		return sendCoinRequest, errors.New("Invalid request body.\n" + err.Error())
	}
	return sendCoinRequest, nil
}

func sendCoinsToUser(fromUser, toUser entities.User, amount uint) (statusCode int, err error) {
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
	routeVariables := getRouteVariables(r)
	merchType := routeVariables["item"]

	merch, err := database.GetMerchByType(merchType)
	if err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	user, err := database.GetUserByUsername(getUsernameFromHeader(r))
	if err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	statusCode, err := buyMerch(user, merch)
	if err != nil {
		respondError(w, statusCode, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func getRouteVariables(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func buyMerch(user entities.User, merch entities.Merch) (statusCode int, err error) {
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
	authRequest, err := getAuthRequest(r)
	if err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	if err = hasUsernameAndPassword(authRequest); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	user, statusCode, err := validateUserOrCreateNewOne(authRequest)
	if err != nil {
		respondError(w, statusCode, err.Error())
	}

	tokenJWT, err := jwtutil.GenerateTokenFor(user)
	if err != nil {
		respondInternalServerError(w, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, model.AuthResponse{Token: tokenJWT})
}

func getAuthRequest(r *http.Request) (authRequest model.AuthRequest, err error) {
	err = json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		return authRequest, errors.New("Invalid request body.\n" + err.Error())
	}
	return authRequest, nil
}

func hasUsernameAndPassword(authRequest model.AuthRequest) error {
	if isEmptyString(authRequest.Username) || isEmptyString(authRequest.Password) {
		return errors.New("username or password is empty")
	}
	return nil
}

func validateUserOrCreateNewOne(authRequest model.AuthRequest) (user entities.User, statusCode int, err error) {
	user, err = isThereUserWithSuchUsername(authRequest.Username)
	if err == nil {
		if err := user.CheckPassword(authRequest.Password); err != nil {
			return user, http.StatusUnauthorized, err
		}
	} else {
		newUserCredentials := entities.User{Username: authRequest.Username, Password: authRequest.Password}
		user, err = createNewUser(newUserCredentials)
		if err != nil {
			return user, http.StatusInternalServerError, err
		}
	}

	return user, http.StatusOK, nil
}

func isThereUserWithSuchUsername(username string) (entities.User, error) {
	return database.GetUserByUsername(username)
}

func createNewUser(newUser entities.User) (entities.User, error) {
	err := database.Instance.Create(&newUser).Error
	if err != nil {
		return newUser, errors.New("Couldn't create a new user.\n" + err.Error())
	}
	return newUser, nil
}

////

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	respondNotFound(w, "Wrong path or this function isn't implemented yet")
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	respondMethodNotAllowed(w, http.StatusText(http.StatusMethodNotAllowed))
}
