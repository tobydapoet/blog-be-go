package models

type Client struct {
	ID            uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID     uint    `gorm:"unique;not null" json:"account_id"`
	Bio           string  `gorm:"type:varchar(2000)" json:"bio"`
	Account       Account `gorm:"foreignKey:AccountID" json:"account"`
	IsDeleted     bool    `gorm:"default:false" json:"isDeleted"`
	LinkInstagram string  `gorm:"varchar(100)" json:"link_instagram"`
	LinkFacebook  string  `gorm:"varchar(100)" json:"link_facebook"`
	LinkWebsite   string  `gorm:"varchar(100)" json:"link_website"`

	Comments   []Comment  `gorm:"foreignKey:ClientID" json:"-"`
	Blog       []Blog     `gorm:"many2many:favourites" json:"-"`
	Blogs      []Blog     `gorm:"foreignKey:ClientID" json:"-"`
	Activity   []Activity `gorm:"foreignKey:ClientID" json:"-"`
	Followings []Client   `gorm:"many2many:followings;joinForeignKey:ClientID;joinReferences:FollowingID"`
	Followers  []Client   `gorm:"many2many:followings;joinForeignKey:FollowingID;joinReferences:ClientID"`
}
