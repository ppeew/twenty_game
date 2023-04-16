package api

import (
	"context"
	"encoding/json"
	"game_web/global"
	"game_web/global/response"
	"game_web/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
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

func ModelToResponse(room models.Room) response.RoomResponse {
	roomResponse := &response.RoomResponse{
		RoomID:     room.RoomID,
		RoomNumber: room.PeopleNumber,
		GameCount:  room.GameCount,
		RoomOwner:  room.RoomOwner,
	}
	return *roomResponse
}

func SayHello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "hello",
	})
}

// 获取所有等待的房间
func GetRoomList(ctx *gin.Context) {

	keys := global.RedisDB.Keys(context.Background(), "*")
	result, err := keys.Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	var data []response.RoomResponse
	for _, key := range result {
		res := global.RedisDB.Get(context.Background(), key)
		if res.Err() != nil {
			continue
		}
		room := models.Room{}
		_ = json.Unmarshal([]byte(res.Val()), &room)
		data = append(data, ModelToResponse(room))
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// 创建房间
func CreateRoom(ctx *gin.Context) {
	//获取jwt
	claims, _ := ctx.Get("claims")
	userID := claims.(*models.CustomClaims).ID

	roomIDStr := ctx.Query("room_id")
	roomID, _ := strconv.Atoi(roomIDStr)
	peopleNumber, _ := strconv.Atoi(ctx.Query("people_number"))
	gameCount, _ := strconv.Atoi(ctx.Query("game_count"))
	roomInfo := models.Room{
		RoomID:       roomID,
		PeopleNumber: peopleNumber,
		GameCount:    gameCount,
		RoomOwner:    int(userID),
	}
	record, _ := json.Marshal(roomInfo)
	//先不考虑并发问题
	cmd := global.RedisDB.SetNX(context.Background(), roomIDStr, record, 0)
	if !cmd.Val() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "已经存在该房间",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建房间成功",
		"data": ModelToResponse(roomInfo),
	})
}

// 删除房间（仅房主）
func DropRoom(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*models.CustomClaims).ID
	roomIDStr := ctx.Query("room_id")
	//先查询房间是否存在
	cmd := global.RedisDB.Get(context.Background(), roomIDStr)
	result, err := cmd.Result()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间号不存在",
		})
	}
	room := models.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	if int(userID) != room.RoomOwner {
		//不是房主
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "非房主，不可删除房间",
		})
		return
	}

	del := global.RedisDB.Del(context.Background(), roomIDStr)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除房间成功",
		"err": del.Err(),
	})
}

// 更新房间的房主或者游戏配置(仅房主)
func UpdateRoom(ctx *gin.Context) {
	//获取jwt
	claims, _ := ctx.Get("claims")
	userID := claims.(*models.CustomClaims).ID

	roomIDStr := ctx.Query("room_id")
	peopleNumber, _ := strconv.Atoi(ctx.DefaultQuery("people_number", "0"))
	gameCount, _ := strconv.Atoi(ctx.DefaultQuery("game_count", "0"))
	owner, _ := strconv.Atoi(ctx.DefaultQuery("owner", "0"))

	//先查询房间是否存在
	cmd := global.RedisDB.Get(context.Background(), roomIDStr)
	result, err := cmd.Result()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间号不存在",
		})
		return
	}
	room := models.Room{}
	_ = json.Unmarshal([]byte(result), &room)
	if int(userID) != room.RoomOwner {
		//不是房主
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "非房主，不可更改房间",
		})
		return
	}
	if peopleNumber != 0 {
		room.PeopleNumber = peopleNumber
	}
	if gameCount != 0 {
		room.GameCount = gameCount
	}
	if owner != 0 {
		room.RoomOwner = owner
	}
	record, _ := json.Marshal(room)

	res := global.RedisDB.Set(context.Background(), roomIDStr, record, 0)
	if res.Err() != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": res.Err(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新房间成功",
		"data": ModelToResponse(room),
	})
}
