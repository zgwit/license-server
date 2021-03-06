package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"license-server/model"
	"net/http"
	"reflect"
)

func mustLogin(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if user := session.Get("user"); user != nil {
		ctx.Set("user", user)
		ctx.Next()
	} else {
		//TODO 检查OAuth2返回的code，进一步获取用户信息，放置到session中

		ctx.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "Unauthorized"})
		ctx.Abort()
	}
}

func RegisterRoutes(app *gin.RouterGroup) {
	//错误恢复，并返回至前端
	app.Use(func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//runtime.Stack()
				//debug.Stack()
				switch err.(type) {
				case error:
					replyError(ctx, err.(error))
				case string:
					replyFail(ctx, err.(string))
				default:
					ctx.JSON(http.StatusOK, gin.H{"error": err})
				}
			}
		}()
		ctx.Next()

		//TODO 内容如果为空，返回404

	})

	app.POST("/login", login)

	app.GET("/active/:code", activeInfo)
	app.POST("/active", curdApiCreate(reflect.TypeOf(model.License{}), licenseBeforeCreate, nil))

	//检查 session，必须登录
	app.Use(mustLogin)

	app.GET("/logout", logout)
	app.POST("/password", password)

	//用户接口
	modelUser := reflect.TypeOf(model.User{})
	app.GET("/user/me", userMe)
	app.POST("/user/list", curdApiList(modelUser))
	app.POST("/user/create", curdApiCreate(modelUser, nil, nil))
	app.GET("/user/:id", parseParamId, curdApiGet(modelUser))
	app.POST("/user/:id", parseParamId, curdApiModify(modelUser, []string{"username", "nickname", "disabled"}, nil, nil))
	app.GET("/user/:id/delete", parseParamId, curdApiDelete(modelUser, nil, nil))
	app.GET("/user/:id/password", parseParamId, userPassword)
	app.GET("/user/:id/enable", parseParamId, curdApiDisable(modelUser, false, nil, nil))
	app.GET("/user/:id/disable", parseParamId, curdApiDisable(modelUser, true, nil, nil))

	//产品接口
	modelProduct := reflect.TypeOf(model.Product{})
	app.POST("/product/list", curdApiList(modelProduct))
	app.POST("/product/create", curdApiCreate(modelProduct, productBeforeCreate, nil))
	app.GET("/product/:id", parseParamId, curdApiGet(modelProduct))
	app.POST("/product/:id", parseParamId, curdApiModify(modelProduct,
		[]string{"code", "name", "uuid", "sn", "cpu", "mac", "disabled"}, nil, nil))
	app.GET("/product/:id/delete", parseParamId, curdApiDelete(modelProduct, nil, nil))
	app.GET("/product/:id/enable", parseParamId, curdApiDisable(modelProduct, false, nil, nil))
	app.GET("/product/:id/disable", parseParamId, curdApiDisable(modelProduct, true, nil, nil))

	//周期接口
	modelTerm := reflect.TypeOf(model.Term{})
	app.POST("/term/list", curdApiList(modelTerm))
	app.POST("/term/create", curdApiCreate(modelTerm, nil, nil))
	app.GET("/term/:id", parseParamId, curdApiGet(modelTerm))
	app.POST("/term/:id", parseParamId, curdApiModify(modelTerm,
		[]string{"name", "description", "year", "month", "day", "price"}, nil, nil))
	app.GET("/term/:id/delete", parseParamId, curdApiDelete(modelTerm, nil, nil))

	//项目接口
	modelLicense := reflect.TypeOf(model.License{})
	app.POST("/license/list", curdApiList(modelLicense))
	app.POST("/license/create", curdApiCreate(modelLicense, licenseBeforeCreate, nil))
	app.GET("/license/:id", parseParamId, curdApiGet(modelLicense))
	app.POST("/license/:id", parseParamId, curdApiModify(modelLicense,
		[]string{"user", "email", "cellphone", "organization"}, nil, nil))
	app.GET("/license/:id/delete", parseParamId, curdApiDelete(modelLicense, nil, nil))

	//TODO 报接口错误（以下代码不生效，路由好像不是树形处理）
	app.Use(func(ctx *gin.Context) {
		replyFail(ctx, "Not found")
		ctx.Abort()
	})
}

func replyList(ctx *gin.Context, data interface{}, total int64) {
	ctx.JSON(http.StatusOK, gin.H{"data": data, "total": total})
}

func replyOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func replyFail(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusOK, gin.H{"error": err})
}

func replyError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
}

func nop(ctx *gin.Context) {
	ctx.String(http.StatusForbidden, "Unsupported")
}
