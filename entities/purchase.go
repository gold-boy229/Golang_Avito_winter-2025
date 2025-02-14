package entities

type Purchase struct {
	Id      uint `gorm:"primaryKey" json:"id"`
	UserId  uint `json:"userId"`
	MerchId uint `json:"merchId"`
	// Quantity uint  `gorm:"default:1" json:"quantity"`
	Owner User  `gorm:"foreignKey:UserId"`
	Merch Merch `gorm:"foreignKey:MerchID"`
}
