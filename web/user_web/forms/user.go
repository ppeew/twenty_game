package forms

type RegisterForm struct {
	Nickname string `form:"nickname" binding:"required,min=1,max=30"`
	Gender   string `form:"gender"`
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
