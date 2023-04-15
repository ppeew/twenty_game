package api

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"user_web/forms"
	"user_web/global"
	"user_web/global/response"
	"user_web/middlewares"
	"user_web/models"
	"user_web/proto"

	"github.com/DanPlayer/randomname"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func ToBool(s string) bool {
	if s == "true" {
		return true
	}
	return false
}

func GetRandomNickName(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"nickname": randomname.GenerateName(),
	})
}

// 查询用户信息,通过id查询
func GetUserInfo(ctx *gin.Context) {
	queryID, _ := strconv.Atoi(ctx.Query("id"))
	info, err := global.UserSrvClient.GetUserByID(context.Background(), &proto.UserIDInfo{Id: uint32(queryID)})
	if err != nil {
		GrpcErrorToHttp(err, ctx)
		return
	}
	resp := response.UserResponse{
		Id:       info.Id,
		Nickname: info.Nickname,
		Gender:   info.Gender,
		Username: info.UserName,
	}
	ctx.JSON(http.StatusOK, resp)
}

func UserRegister(ctx *gin.Context) {
	register := forms.RegisterForm{}
	if err := ctx.ShouldBind(&register); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}

	var info *proto.UserInfoResponse
	var err error
	//保证一定能够生成，多次尝试
	for true {
		info, err = global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			Nickname: register.Nickname,
			Gender:   ToBool(register.Gender),
			UserName: randomname.GenerateName(),
			Password: "123456",
		})
		if err == nil {
			break
		}
		code, _ := status.FromError(err)
		if code.Code() != codes.AlreadyExists {
			//服务错误
			//GrpcErrorToHttp(err, ctx)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err,
			})
			return
		}
		//数据库存在，持续循环直到可以添加
	}

	//返回token，实际上默认已经完成登录过程了
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID: info.Id,
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
			"msg": err.Error(),
			"err": 0,
		})
		return
	}
	info, err := global.UserSrvClient.GetUserByUsername(context.Background(), &proto.UserNameInfo{UserName: login.Username})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户不存在",
			"err": 1,
		})
		return
	}
	verify, err := global.UserSrvClient.CheckPassword(context.Background(), &proto.CheckPasswordInfo{
		Password:       login.Password,
		EncodePassword: info.Password,
	})
	if err != nil {
		GrpcErrorToHttp(err, ctx)
	}
	if !verify.Success {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "密码错误",
			"err": 2,
		})
		return
	}
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID: info.Id,
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
			"msg": err.Error(),
		})
		return
	}

	//必须先查询是否有username=form.Username
	_, err := global.UserSrvClient.GetUserByUsername(context.Background(), &proto.UserNameInfo{UserName: form.Username})
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户名已经被使用",
		})
		return
	}

	_, err = global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
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
		"msg": "更改信息成功",
	})
}
