package models

type User struct {
	BaseModel
	UserName    string `gorm:"type:string;not null;unique;size:60"`
	PhoneNumber string `gorm:"type:string;not null;unique;size:60"`
	Password    string `gorm:"type:string;not null"`
	UserRoles   []UserRole
	Verified    bool `gorm:"type:boolean;default:false"`
	WishList    []UserWishList
}

type Role struct {
	BaseModel
	Name      string `gorm:"type:string;not null;unique;size:60"`
	UserRoles []UserRole
}

type UserRole struct {
	UserId int
	RoleId int
	Role   Role `gorm:"foreignKey:RoleId"`
	User   User `gorm:"foreignKey:UserId"`
}

// TODO: add preSave hook
