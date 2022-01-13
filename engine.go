package leowebgin

import (
	"fmt"
	"github.com/cqlsy/leolog"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"net/http"
	"time"
)

type WebGin struct {
	Gin *gin.Engine
}

//	c *gin.Context
// 1: filePath := c.Query("filePath")
// 2: filePath = c.PostForm("filePath")
// 3: data := make(map[string]interface{})
//    c.ShouldBind(&data)
func New(runNode string) *WebGin {
	//gin.DefaultWriter = ioutil.Discard
	if runNode == "pro" {
		gin.SetMode(gin.ReleaseMode)
	}
	web := new(WebGin)
	//web.Gin = gin.Default()
	web.Gin = gin.New()
	// set common middle func
	web.Gin.Use(cors(), gzip.Gzip(gzip.DefaultCompression), middleLog(), gin.Recovery())
	return web
}

//  ip:0.0.0.0  port:8080
func (w *WebGin) StartListen(ip interface{}, port interface{}) {
	str := fmt.Sprintf("%v:%v", ip, port)
	leolog.Print("Web Listener On:" + str)
	err := w.Gin.Run(str)
	if err != nil {
		panic(fmt.Sprintf("Listener %s error: %s", str, err.Error()))
	}
}

// packr.New("static", "./static")
//  http.Dir("./static")
func (w *WebGin) AddStaticPath(routePath, dirPath string) {
	w.Gin.StaticFS(routePath, http.Dir(dirPath))
}

// (go build -o xxx) change to (packr build -o xxx)
func (w *WebGin) AddStaticPackrPath(path, packrName, dirPath string) {
	if path == "/" {
		fmt.Println("you can't set GET Request when you set static path with '/'. " +
			"Are you sure do this? ")
	}
	w.Gin.StaticFS(path, packr.New(packrName, dirPath))
	//w.Gin.StaticFS(path, fs)
}

func cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Token,"+
			"Referer,X-Requested-With")
		context.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			//context.AbortWithStatus(http.StatusNoContent)
			context.JSON(http.StatusOK, "options request")
		}
		context.Next()
	}
}

func middleLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.Request.Host
		leolog.LogInfoDefault(fmt.Sprintf("| code: %3d | time: %13v | clientIP: %15s | reqMethod: %s | reqUri: %s |",
			statusCode, latencyTime, clientIP, reqMethod, reqUri))
	}
}

func GetReqParams(c *gin.Context, name string) string {
	params := c.Query(name)
	if params == "" {
		params = c.PostForm(name)
	}
	return params
}

//     router.GET("/user/:id", func(c *gin.Context) {
//         // a GET request to /user/john
//         id := c.Param("id") // id == "john"
//     })
func GetReqParamsFromPath(c *gin.Context, name string) string {
	params := c.Param(name)
	return params
}

func GetReqParamsFromBody(c *gin.Context) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := c.ShouldBind(&data)
	return data, err
}
