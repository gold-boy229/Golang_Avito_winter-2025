package entities

type Transaction struct {
	Id         uint `gorm:"primaryKey" json:"id"`
	FromUserId uint `json:"fromUserId"`
	ToUserId   uint `json:"toUserId"`
	Amount     uint `json:"amount"`
	FromUser   User `gorm:"foreignKey:FromUserID"`
	ToUser     User `gorm:"foreignKey:ToUserID"`
}
