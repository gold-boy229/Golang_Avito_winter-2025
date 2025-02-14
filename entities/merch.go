package entities

type Merch struct {
	Id   uint   `json:"id" gorm:"primaryKey"`
	Type string `json:"type"`
	Cost uint   `json:"cost"`
}
