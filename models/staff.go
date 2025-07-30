package models

type Staff struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID uint    `gorm:"unique;not null" json:"account_id"`
	Phone     string  `gorm:"type:varchar(20);unique;not null" json:"phone"`
	Account   Account `gorm:"foreignKey:AccountID" json:"account"`
	IsDeleted bool    `gorm:"default:false" json:"isDeleted"`
}
