package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"user_web/forms"
	"user_web/global"
	"user_web/global/response"
	"user_web/middlewares"
	"user_web/models"
	"user_web/proto/user"

	"go.uber.org/zap"

	"github.com/DanPlayer/randomname"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetRandomNickName(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"nickname": randomname.GenerateName(),
	})
}

func GetRandomUsername(ctx *gin.Context) {
	username := fmt.Sprintf("%05v", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000))
	ctx.JSON(http.StatusOK, gin.H{
		"username": username,
	})
}

// 查询用户信息,通过id查询
func GetUserInfo(ctx *gin.Context) {
	queryID, _ := strconv.Atoi(ctx.Query("id"))
	info, err := global.UserSrvClient.GetUserByID(context.Background(), &user.UserIDInfo{Id: uint32(queryID)})
	if err != nil {
		GrpcErrorToHttp(err, ctx)
		return
	}
	resp := response.UserResponse{
		Id:       info.Id,
		Nickname: info.Nickname,
		Gender:   info.Gender,
		Username: info.UserName,
		Image:    info.Image,
		State:    info.State,
	}
	ctx.JSON(http.StatusOK, resp)
}

func UserRegister(ctx *gin.Context) {
	register := forms.RegisterForm{}
	if err := ctx.ShouldBind(&register); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	var info *user.UserInfoResponse
	var err error
	info, err = global.UserSrvClient.CreateUser(context.Background(), &user.CreateUserInfo{
		Nickname: register.Nickname,
		Gender:   ToBool(register.Gender),
		UserName: register.Username,
		Password: register.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	//返回token，实际上默认已经完成登录过程了
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID: info.Id,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),              //签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*5, //过期时间.五天
			Issuer:    "ppeew",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": "生成token失败",
		})
		return
	}
	resp := response.UserResponse{
		Id:       info.Id,
		Nickname: info.Nickname,
		Gender:   info.Gender,
		Username: info.UserName,
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000, //返回是毫秒的
		"data":       resp,
	})
}

// 用户登录
func UserLogin(ctx *gin.Context) {
	login := forms.LoginForm{}
	if err := ctx.ShouldBind(&login); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	info, err := global.UserSrvClient.GetUserByUsername(context.Background(), &user.UserNameInfo{UserName: login.Username})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "用户没注册",
		})
		return
	}
	verify, err := global.UserSrvClient.CheckPassword(context.Background(), &user.CheckPasswordInfo{
		Password:       login.Password,
		EncodePassword: info.Password,
	})
	if err != nil {
		GrpcErrorToHttp(err, ctx)
	}
	if !verify.Success {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "密码错误",
		})
		return
	}
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID: info.Id,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),              //签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*5, //过期时间.五天
			Issuer:    "ppeew",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": "生成token失败",
		})
		return
	}
	resp := response.UserResponse{
		Id:       info.Id,
		Nickname: info.Nickname,
		Gender:   info.Gender,
		Username: info.UserName,
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000, //返回是毫秒的
		"data":       resp,
	})

}

// 当前用户修改信息
func UserUpdate(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	id := currentUser.ID

	form := forms.ModifyForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	//必须先查询是否有username=form.Username
	_, err := global.UserSrvClient.GetUserByUsername(context.Background(), &user.UserNameInfo{UserName: form.Username})
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": "用户名已经被使用",
		})
		return
	}
	_, err = global.UserSrvClient.UpdateUser(context.Background(), &user.UpdateUserInfo{
		Id:       id,
		Nickname: form.Nickname,
		Gender:   ToBool(form.Gender),
		UserName: form.Username,
		Password: form.Password,
	})
	if err != nil {
		GrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": "更改信息成功",
	})
}

// 获得用户的状态
func SelectUserState(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*models.CustomClaims).ID
	state, err := global.UserSrvClient.GetUserState(context.Background(), &user.UserIDInfo{Id: userID})
	if err != nil {
		zap.S().Warnf("[SelectUserState]:%s", err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": state.State,
		"err":  "",
	})
}

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

func ToBool(s string) bool {
	if s == "true" {
		return true
	}
	return false
}
