package models

type UpdateUserModel struct {
	UserID    string  `json:"userId" validate:"required,uuid"`
	FirstName *string `json:"firstName,omitempty" validate:"omitempty,min=3,max=32"`
	LastName  *string `json:"lastName,omitempty" validate:"omitempty,min=3,max=32"`
	Email     *string `json:"email" validate:"email,max=255,min=6"`
	Phone     *string `json:"phone" validate:"omitempty,min=8,max=15"`
}

type RegisterModel struct {
	Username  string         `json:"username" validate:"required,min=3,max=32,lowercase"`
	Password  string         `json:"password" validate:"required,min=6,max=32"`
	FirstName string         `json:"firstName,omitempty" validate:"omitempty,min=3,max=32"`
	LastName  string         `json:"lastName,omitempty" validate:"omitempty,min=3,max=32"`
	Phone     string         `json:"phone" validate:"required,min=10,max=15"`
	Email     string         `json:"email" validate:"required,email,max=255,min=6"`
	Address   *CreateAddress `json:"address" validate:"omitempty,required"`
}

type LoginModel struct {
	Username *string `form:"username" validate:"omitempty,max=32"`
	Email    *string `form:"email" validate:"omitempty,email,max=255"`
	Password string  `form:"password" validate:"required,min=6,max=32"`
}

type VerifyEmailQuery struct {
	VerifyCode string `form:"verifyCode" validate:"required,min=1"`
}
