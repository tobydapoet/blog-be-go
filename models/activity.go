package models

import (
	"time"

	"gorm.io/datatypes"
)

type Activity struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID  uint           `gorm:"not null" json:"clientId"`
	Content   string         `gorm:"type:varchar(2000);not null" json:"content"`
	Images    datatypes.JSON `gorm:"type:json" json:"images"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	IsDeleted bool           `gorm:"default:false" json:"isDeleted"`

	Client Client `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE" json:"client"`
}
