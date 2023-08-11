package dto

type CommentListReq struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type CommentReq struct {
	ID      int    `form:"id" uri:"id"`
	Content string `form:"content" uri:"content"`
}
