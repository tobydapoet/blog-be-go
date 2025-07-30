package models

import "time"

type BlackListToken struct {
	ID        uint      `gorm:"primaryKey,autoIncrement"`
	Token     string    `gorm:"size:255;uniqueIndex" json:"token"`
	Type      string    `gorm:"type:enum('acess','refresh');default:'refresh';not null" json:"type"`
	AccountID uint      `gorm:"not null" json:"account_id"`
	ExpiresAt time.Time `gorm:"type: datetime" json:"expiresAt"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Account Account `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE" json:"account"`
}
