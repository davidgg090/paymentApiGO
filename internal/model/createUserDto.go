package model

type createUserDTO struct {
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;unique" json:"password"`
}
