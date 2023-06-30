package forms

type RegisterForm struct {
	Nickname string `form:"nickname" binding:"required,min=1,max=30"`
	Gender   string `form:"gender" binding:"required"`
	Username string `form:"username" binding:"required,min=1,max=5"`
	Password string `form:"password" binding:"required,max=20"`
}

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type ModifyForm struct {
	Nickname string `form:"nickname"`
	Gender   string `form:"gender"`
	Username string `form:"username"`
	Password string `form:"password"`
}
