package models

type Category struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"varchar(500);not null" json:"name"`
	IsDeleted bool   `gorm:"default:false" json:"isDeleted"`

	Blogs []Blog `gorm:"many2many:blog_categories" json:"blogs"`
}
