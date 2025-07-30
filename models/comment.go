package models

import "time"

type Comment struct {
	ID               uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CommenttableID   uint   `gorm:"not null" json:"commentTableId"`
	CommenttableType string `gorm:"enum('blog','activity',comment)" json:"commentTableType"`

	ClientID  uint      `gorm:"not null" json:"client_id"`
	Content   string    `gorm:"type:varchar(1000);not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	IsDeleted bool      `gorm:"default:false" json:"isDeleted"`

	Client Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

	Commenttable any `gorm:"-" json:"commentTable"`
}
