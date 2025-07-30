package models

type Following struct {
	ID          uint `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID    uint `gorm:"not null" json:"client_id"`
	FollowingID uint `gorm:"not null" json:"following_id"`

	Client       Client `gorm:"foreignKey:ClientID" json:"client"`
	FollowedUser Client `gorm:"foreignKey:FollowingID" json:"followed_user"`
}
