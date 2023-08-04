package tests

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestBind(t *testing.T) {
	router := gin.Default()
	//测试
	type P1 struct {
		ID int `uri:"id"`
	}
	router.GET("/bindUri/:id", func(ctx *gin.Context) {
		var d P1
		println(ctx.Param("id"))
		if err := ctx.ShouldBindUri(&d); err != nil { //1.当不写shouldBIndURi时，是查询不到id的
			ctx.Status(400)
		} else {
			ctx.Status(200)
		}
	})

	type P2 struct {
		ID int `form:"id"`
	}
	router.GET("/bindQuery", func(ctx *gin.Context) {
		var d P2
		println(ctx.Query("id"))
		if err := ctx.ShouldBindQuery(&d); err != nil {
			ctx.Status(400)
		} else {
			ctx.Status(200)
		}
	})

	type P3 struct {
		ID int `form:"id"`
	}
	router.POST("/bindForm", func(ctx *gin.Context) {
		var d P3
		if err := ctx.ShouldBind(&d); err != nil {
			//必须要post请求，shouldBind绑定form表单数据
			// If `GET`, only `Form` binding engine (`query`) used.
			// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
			ctx.Status(400)
		} else {
			ctx.Status(200)
		}
	})
	router.Run(":10000")
}
