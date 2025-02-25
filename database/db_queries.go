package database

import (
	"MerchShop/entities"
	"fmt"
)

func GetMerchByType(merchType string) (merch entities.Merch, err error) {
	result := Instance.Where("type = ?", merchType).First(&merch)
	if result.Error != nil {
		err = fmt.Errorf("there is no merch with type %s", merchType)
	}
	return merch, err
}
