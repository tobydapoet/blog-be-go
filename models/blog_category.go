package models

type BlogCategory struct {
	ID         uint     `gorm:"primaryKey;autoIncrement"`
	BlogID     uint     `json:"blogId"`
	CategoryID uint     `json:"categoryId"`
	Blog       Blog     `gorm:"foreignKey:BlogID" json:"blog"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`
}
