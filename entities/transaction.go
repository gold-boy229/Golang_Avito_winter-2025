package entities

type Transaction struct {
	Id         uint `json:"id" gorm:"primaryKey"`
	FromUserId uint `json:"fromUserId"`
	ToUserId   uint `json:"toUserId"`
	Amount     uint `json:"amount"`
	FromUser   User `gorm:"foreignKey:FromUserId;references:Id"`
	ToUser     User `gorm:"foreignKey:ToUserId;references:Id"`
}
