package api

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"time"
	"user_web/global"
	"user_web/global/response"
	"user_web/middlewares"
	"user_web/models"
	"user_web/proto"
)

// GRPC错误转HTTP
func GrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Message(),
				})
			}
		}
	}
}

// 查询用户信息,通过id查询
func GetUserInfo(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户:%d", currentUser.ID)

	queryID, _ := strconv.Atoi(ctx.Query("id"))
	info, err := global.UserSrvClient.GetUserByID(context.Background(), &proto.UserIDInfo{Id: uint32(queryID)})
	if err != nil {
		GrpcErrorToHttp(err, ctx)
	}
	resp := response.UserResponse{
		Id:     info.Id,
		Name:   info.Name,
		Gender: info.Gender,
	}
	ctx.JSON(http.StatusOK, resp)
}

// 用户注册,提交code
func UserRegister(ctx *gin.Context) {
	//通过code获取openID 									TODO TODO
	openID := "1111"

	info, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Name:   ctx.Query("name"),
		OpenID: openID,
		Gender: ctx.Query("gender"),
	})

	if err != nil {
		code, _ := status.FromError(err)
		if code.Code() != codes.AlreadyExists {
			//服务错误
			GrpcErrorToHttp(err, ctx)
		}
	}

	//服务正常运行(忽略是否真正添加了用户，只要用户在表里，返回客户端token即可)
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:   uint32(info.Id),
		Name: info.Name,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),            //签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24, //过期时间.一天
			Issuer:    "ppeew",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":        "成功",
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000, //返回是毫秒的
	})
}

// 用户登录(用token印证)
func UserLogin(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("id为%d用户访问登录接口", currentUser.ID)
	// 由于用户已经token印证过了，确认该用户身份，直接返回成功
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// 当前用户修改信息
func UserUpdate(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	id := currentUser.ID
	name := ctx.Query("name")
	gender := ctx.Query("gender")

	_, err := global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Name:   name,
		Gender: gender,
		Id:     id,
	})
	if err != nil {
		GrpcErrorToHttp(err, ctx)
	}
	resp := response.UserResponse{
		Id:     id,
		Name:   name,
		Gender: gender,
	}
	ctx.JSON(http.StatusOK, resp)
}
