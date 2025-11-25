package models

type UpdateUserModel struct {
	UserID    string  `json:"userId" binding:"required,uuid"`
	FirstName *string `json:"firstName,omitempty" binding:"omitempty,min=3,max=32"`
	LastName  *string `json:"lastName,omitempty" binding:"omitempty,min=3,max=32"`
	Email     *string `json:"email" binding:"email,max=255,min=6"`
	Phone     *string `json:"phone" binding:"omitempty,min=8,max=15"`
}

type RegisterModel struct {
	Username  string         `json:"username" binding:"required,min=3,max=32,lowercase"`
	Password  string         `json:"password" binding:"required,min=6,max=32"`
	FirstName string         `json:"firstName,omitempty" binding:"omitempty,min=3,max=32"`
	LastName  string         `json:"lastName,omitempty" binding:"omitempty,min=3,max=32"`
	Phone     string         `json:"phone" binding:"required,min=10,max=15"`
	Email     string         `json:"email" binding:"required,email,max=255,min=6"`
	Address   *CreateAddress `json:"address" binding:"omitempty,required"`
}

type LoginModel struct {
	Username *string `form:"username" binding:"omitempty,max=32"`
	Email    *string `form:"email" binding:"omitempty,email,max=255"`
	Password string  `form:"password" binding:"required,min=6,max=32"`
}

type VerifyEmailQuery struct {
	VerifyCode string `form:"verifyCode" binding:"required,min=1"`
}
