package apis

import (
	"github.com/gin-gonic/gin"
	"hall_web/model"
	"hall_web/service/dao"
	"hall_web/service/domains"
	"hall_web/service/dto"
	"net/http"
	"time"
)

func CommentList(ctx *gin.Context) {
	req := dto.CommentListReq{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	commentDao := &dao.CommentDao{}
	res, err := commentDao.CommentList(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "获取数据错误")
		return
	}

	if len(res) == 0 {
		ctx.JSON(http.StatusOK, "无留言")
	}
	ctx.JSON(http.StatusOK, res)
}

func AddComment(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID

	content := ctx.PostForm("content")
	if content == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	c := domains.Comments{
		UserID:  int(userID),
		Time:    time.Now().Format("2006-01-02 15:04"),
		Content: content,
	}
	commentDao := dao.CommentDao{}
	err := commentDao.AddComment(&c)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, "添加评论成功")
}

func UpdateComment(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID

	req := dto.CommentReq{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	c := domains.Comments{
		ID:      req.ID,
		UserID:  int(userID),
		Time:    time.Now().Format("2006-01-02 15:04"),
		Content: req.Content,
	}
	commentDao := dao.CommentDao{}
	err = commentDao.UpdateComment(&c)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, "更新评论成功")
}

func DelComment(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID

	req := dto.CommentReq{}
	err := ctx.BindUri(&req)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	c := domains.Comments{
		ID:     req.ID,
		UserID: int(userID),
	}
	commentDao := dao.CommentDao{}
	err = commentDao.DelComment(&c)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, "删除评论成功")
}

// 暂时不做，表需要增加点赞玩家id列表的字段，以后完善
func LikeComment(ctx *gin.Context) {
	req, exist := ctx.Get("user_id")
	if !exist {
		ctx.Status(http.StatusBadRequest)
		return
	}

	c := domains.Comments{
		UserID: req.(int),
	}
	commentDao := dao.CommentDao{}
	err := commentDao.LikeComment(&c)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, "点赞评论成功")
}
