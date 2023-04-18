package api

import (
	"context"
	"encoding/json"
	"fmt"
	"game_web/global"
	"game_web/model"
	"game_web/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

func SayHello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "hello",
	})
}

// 获取所有的房间
func GetRoomList(ctx *gin.Context) {
	keys := global.RedisDB.Keys(context.Background(), "*")
	if keys.Err() != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": keys.Err(),
		})
		return
	}
	var resp []model.Room
	for _, key := range keys.Val() {
		var room model.Room
		get := global.RedisDB.Get(context.Background(), key)
		result, err := get.Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
			return
		}
		_ = json.Unmarshal([]byte(result), &room)
		resp = append(resp, room)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// 创建房间
func CreateRoom(ctx *gin.Context) {
	//获取jwt
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID

	roomIDStr := ctx.Query("room_id")
	roomID, _ := strconv.Atoi(roomIDStr)
	maxUserNumber, _ := strconv.Atoi(ctx.Query("max_user_number"))
	gameCount, _ := strconv.Atoi(ctx.Query("game_count"))

	room := model.Room{RoomID: roomID}
	if _, err := room.Select(room); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间已经存在",
		})
		return
	}

	room = model.Room{
		RoomID:        roomID,
		MaxUserNumber: maxUserNumber,
		GameCount:     gameCount,
		UserNumber:    0,
		RoomOwner:     userID,
		RoomWait:      true,
		Users:         nil,
		Publish:       model.NewPublisher(), //创建发布者
	}
	err := room.Create(room)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "创建房间失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "创建房间成功",
	})
}

// 删除房间（仅房主）
func DropRoom(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	//先查询房间是否存在
	room := model.Room{RoomID: roomID}
	retRoom, err := room.Select(room)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": "找不到该房间",
		})
	}

	if userID != retRoom.RoomOwner {
		//不是房主
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "非房主，不可删除房间",
		})
		return
	}
	_ = retRoom.Delete(retRoom)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除房间成功",
	})
}

// 更新房间的房主或者游戏配置(仅房主)
func UpdateRoom(ctx *gin.Context) {
	//获取jwt
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID

	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	maxUserNumber, _ := strconv.Atoi(ctx.DefaultQuery("max_user_number", "0"))
	gameCount, _ := strconv.Atoi(ctx.DefaultQuery("game_count", "0"))
	owner, _ := strconv.Atoi(ctx.DefaultQuery("owner", "0"))
	kicker, _ := strconv.Atoi(ctx.DefaultQuery("kicker", "0"))

	//先查询房间是否存在
	room := model.Room{RoomID: roomID}
	retRoom, err := room.Select(room)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间不存在",
		})
		return
	}
	if userID != retRoom.RoomOwner {
		//不是房主
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "非房主，不可更改房间",
		})
		return
	}

	if maxUserNumber != 0 {
		retRoom.MaxUserNumber = maxUserNumber
	}
	if gameCount != 0 {
		retRoom.GameCount = gameCount
	}
	if owner != 0 {
		retRoom.RoomOwner = owner
	}
	if kicker != 0 && kicker == userID {
		//t房主(不允许这样做)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "不可T自己",
		})
		return
	}
	err = retRoom.Update(retRoom)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "更新房间失败",
		})
		return
	}
	//更新成功，发布更新信息
	marshal, _ := json.Marshal(retRoom)
	retRoom.Publish.Publish(marshal)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新房间成功",
	})
}

// 玩家进入
func UserIntoRoom(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID

	//查找房间是否存在
	room := model.Room{RoomID: roomID}
	retRoom, err := room.Select(room)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	//房间存在，房间当前人数不应该满了或者房间开始了
	if retRoom.UserNumber >= retRoom.MaxUserNumber {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间满了",
		})
		return
	} else if !retRoom.RoomWait {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间已开始游戏",
		})
		return
	}
	//进入房间,建立websocket连接
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "无法连接房间服务器",
		})
		return
	}
	retRoom.Users[userID].WSConn = model.InitWebSocket(conn)
	retRoom.Users[userID].Ready = false
	retRoom.UserNumber++
	err = retRoom.Update(retRoom)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}
	//订阅房间
	room.Publish.AddSubscriber(room.Users[userID].WSConn.Subscriber)

	//因为房间更新，给所有订阅者发送房间信息
	marshal, _ := json.Marshal(retRoom)
	room.Publish.Publish(marshal)
}

// 房间信息
func RoomInfo(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	room := model.Room{RoomID: roomID}
	retRoom, err := room.Select(room)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间不存在",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": retRoom,
	})
}

// 玩家准备状态
func UpdateUserReadyState(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	readyState := ctx.Query("ready_state")
	room := model.Room{RoomID: roomID}
	retRoom, err := room.Select(room)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间不存在",
		})
		return
	}
	if userID == retRoom.RoomOwner {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房主不能准备",
		})
		return
	}
	retRoom.Users[userID].Ready = utils.StringToBool(readyState)
	_ = retRoom.Update(retRoom)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("修改准备状态为：%s", readyState),
	})
	//房间状态改变,发布
	marshal, _ := json.Marshal(retRoom)
	retRoom.Publish.Publish(marshal)
}

// 开始游戏按键接口
func BeginGame(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	claims, _ := ctx.Get("claims")
	userID := claims.(*model.CustomClaims).ID
	//查看房间是否存在
	room := model.Room{RoomID: roomID}
	retRoom, err := room.Select(room)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "房间不存在",
		})
		return
	}
	if userID != retRoom.RoomOwner {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "仅房主才能开始游戏",
		})
		return
	}
	for _, user := range retRoom.Users {
		if user.ID != userID && user.Ready == false {
			//非房主外的用户没准备
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "有玩家没准备",
			})
			return
		}
	}
	//都准备好了，游戏环境，可以进入游戏模块,发布者向所有用户发送游戏开始，TODO
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "可以开始游戏",
	})
}
