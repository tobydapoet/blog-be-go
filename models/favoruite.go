package models

type Favourite struct {
	ID                 uint   `gorm:"primaryKey" json:"id"`
	ClientID           uint   `json:"client_id"`
	FavouritetableID   uint   `gorm:"not null" json:"favouriteTableId"`
	FavouritetableType string `gorm:"type:enum('blog','activity','comment');not null" json:"favouriteTableType"`

	Client Client `gorm:"foreignKey:ClientID" json:"client"`
}
