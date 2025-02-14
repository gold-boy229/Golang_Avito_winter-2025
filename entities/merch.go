package entities

type Merch struct {
	Id   uint   `gorm:"primaryKey" json:"id"`
	Type string `json:"type"`
	Cost uint   `json:"cost"`
}
