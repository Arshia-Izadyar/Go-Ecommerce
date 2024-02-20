package models

type Shipping string

const (
	Express Shipping = "express"
	Plus    Shipping = "plus"
	Normal  Shipping = "normal"
)

type Product struct {
	BaseModel
	Name        string     `gorm:"type:varchar(100);unique;not null"`
	Description string     `gorm:"type:text;"`
	Price       float64    `gorm:"type:decimal(10,2);not null"`
	Images      []string   `gorm:"type:json;"`
	Color       []string   `gorm:"type:json;"`
	InStock     bool       `gorm:"type:boolean;default:true;not null"`
	Quantity    int        `gorm:"type:int;default:1"`
	Shipping    Shipping   `gorm:"type:varchar(20);"`
	Slug        string     `gorm:"type:varchar(100);unique;not null"`
	Categories  []Category `gorm:"many2many:product_categories"`
	Ratings     []Rating   `gorm:"foreignKey:ProductId"`
}

type Category struct {
	BaseModel
	Name    string    `gorm:"type:varchar(100);unique;not null"`
	Slug    string    `gorm:"type:varchar(100);unique;not null"`
	Images  []string  `gorm:"type:json;"`
	Product []Product `gorm:"many2many:product_categories"`
}

type ProductCategory struct {
	ProductId  int `gorm:"primaryKey"`
	CategoryId int `gorm:"primaryKey"`
}

type UserWishList struct {
	BaseModel
	User      User `gorm:"foreignKey:UserId"`
	UserId    int
	Product   Product `gorm:"foreignKey:ProductId"`
	ProductId int
}

type Rating struct {
	BaseModel
	UserId    int     `gorm:"not null"`
	User      User    `gorm:"foreignKey:UserId"`
	Product   Product `gorm:"foreignKey:ProductId"`
	ProductId int     `gorm:"not null;index:user_product_index,unique"`
	Rate      int     `gorm:"type:int"`
	Review    string  `gorm:"type:varchar(1000);null"`
}
