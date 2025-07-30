package models

import "time"

type Blog struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID   uint       `gorm:"not null" json:"clientId"`
	Thumbnail  string     `gorm:"type:text;not null" json:"thumbnail"`
	Title      string     `gorm:"type:varchar(200);not null" json:"title"`
	Content    string     `gorm:"type:varchar(2000);not null" json:"content"`
	Status     string     `gorm:"type:enum('wait approve','approved','denied');default:'wait approve';not null" json:"status"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	Client     Client     `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE" json:"client"`
	Categories []Category `gorm:"many2many:blog_categories" json:"categories"`
	// Accounts   []Account  `gorm:"many2many:favourites" json:"accounts"`
	IsDeleted bool `gorm:"default:false" json:"isDeleted"`
}
