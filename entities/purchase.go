package entities

type Purchase struct {
	Id      uint  `json:"id" gorm:"primaryKey"`
	UserId  uint  `json:"userId"`
	MerchId uint  `json:"merchId"`
	Owner   User  `gorm:"foreignKey:UserId"`
	Merch   Merch `gorm:"foreignKey:MerchId"`
}
