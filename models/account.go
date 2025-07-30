package models

import "time"

type Account struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Password  string    `gorm:"type:varchar(100)" json:"password,omitempty"`
	AvatarURL string    `gorm:"type:text" json:"avatar_url,omitempty"`
	Role      string    `gorm:"type:enum('admin','staff','client');default:'client';not null" json:"role"`
	GoogleID  string    `gorm:"type:varchar(100)" json:"googleId,omitempty"`
	IsDeleted bool      `gorm:"default:false" json:"isDeleted"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Staff  *Staff  `json:"-"`
	Client *Client `json:"-"`
}
