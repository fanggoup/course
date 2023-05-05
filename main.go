package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

// 自定义go中间件
func myHandler() (gin.HandlerFunc){
	return func(ctx *gin.Context) {
		// 通过自定义的中间件，设置的值，在后续处理只要调用了这个中间件的都可以拿到这里的参数
		ctx.Set("usersession","userid-1")
		ctx.Next()//放行
		ctx.Abort()//阻止
	}
}

func main(){
	// 创建一个服务
	ginServer := gin.Default()

	//favicon修改网站图标
	ginServer.Use(favicon.New("./conf/flower.png"))
	//连接数据库代码

	// 加载静态页面
	ginServer.LoadHTMLGlob("./templates/*")

	// 响应一个页面给前端
	ginServer.GET("/index",func(c *gin.Context){
		c.HTML(http.StatusOK,"index.html",gin.H{
			"msg":"这是go后台传输的数据",
		})
	})

	// 接收前端传递过来的参数
	ginServer.GET("/user/info",myHandler(),func(c *gin.Context){

		// 取出中间件中的值
		usersession := c.MustGet("usersession").(string)
		fmt.Println(usersession)

		userid := c.Query("userid")
		username := c.Query("username")
		c.JSON(http.StatusOK,gin.H{
			"userid":userid,
			"username":username,
		})
	})

	// restful传参
	ginServer.GET("/user/info/:username/:userid",func(c *gin.Context){
		userid := c.Param("userid")
		username := c.Param("username")
		c.JSON(http.StatusOK,gin.H{
			"userid":userid,
			"username":username,
		})
	})

	// 前端给后端传JSON
	ginServer.POST("/json",func(c *gin.Context){
		// request.body
		b,_ := c.GetRawData()
		
		var m map[string]interface{}
		_ = json.Unmarshal(b,&m)
		c.JSON(http.StatusOK,m)
	})

	ginServer.POST("/user/add",func(c *gin.Context){
		username := c.PostForm("username")
		password := c.PostForm("password")
		c.JSON(http.StatusOK,gin.H{
			"msg":"ok",
			"username":username,
			"password":password,
		})
	})



	//访问地址，处理我们的请求  request response
	ginServer.GET("/hello",func(context *gin.Context){
		context.JSON(200,gin.H{"msg":"hello"})
	})
	//restful api
	ginServer.POST("/hello")
	ginServer.PUT("/hello")
	ginServer.DELETE("/hello")

	// 路由
	ginServer.GET("/test",func(ctx *gin.Context) {
		// 重定向
		ctx.Redirect(http.StatusMovedPermanently,"https://baidu.com")
	})

	// 404 NoRoute
	ginServer.NoRoute(func(ctx *gin.Context) {
		// why:为什么这里写的地址不是绝对也不是相对地址就行呀，
		ctx.HTML(http.StatusNotFound,"404.html",nil)
	})

	// 路由组
	userGroup := ginServer.Group("/user")
	{
		userGroup.GET("/add")
	}

	// 中间件

	//服务器端口
	ginServer.Run(":8080")
}