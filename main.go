package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	//上传文件测试
	r.POST("/upload", func(c *gin.Context) {
		// single file
		file, _ := c.FormFile("file")
		log.Println(file.Filename)
		dst := "/Users/ouseizan/mygo/src/code.byted.org/baike/myGin/ht.xlsx"
		// Upload the file to specific dst.
		c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// Simple group: v1
	v1 := r.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		//http://localhost:8080/v1/loginHtml?user=wqs&password=565
		v1.POST("/loginHtml", authLoginEndpoint)
		v1.POST("/loginJson", authJsonLoginEndpoint)
	}

	return r
}

func loginEndpoint(c *gin.Context) {
	user := c.Params.ByName("name")
	//注意这两种用法的不同，http://localhost:8080/v1/login?name=wqs  下面这个获取方法才能取到值
	name, _ := c.GetQuery("name")
	user = "登陆者：" + user + name
	value, ok := db[user]
	if ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	}
}

//校验登录
// Example for binding a HTML form (user=manu&password=123)
func authLoginEndpoint(c *gin.Context) {
	var json Login
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if json.User != "manu" || json.Password != "123" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

//在body里写：{"user": "manu", "password": "123"}
func authJsonLoginEndpoint(c *gin.Context) {
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if json.User != "manu" || json.Password != "123" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
