package leowebgin

import (
	"fmt"
	"github.com/cqlsy/leolog"
	"github.com/cqlsy/leotime"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"time"
)

// 新建一个 发动机
func NewEngine(runMode RunMode) *Engine {
	//gin.DefaultWriter = ioutil.Discard
	if runMode == Pro {
		gin.SetMode(gin.ReleaseMode)
	}
	web := new(Engine)
	web.runMode = runMode
	//web.Gin = gin.Default()
	web.engine = gin.New()
	// set common middle func
	web.engine.Use(cors(), gzip.Gzip(gzip.DefaultCompression), web.middleLog(), web.recovery())
	return web
}

//  ip:0.0.0.0  port:8080
func (w *Engine) StartListen(ip interface{}, port interface{}) {
	str := fmt.Sprintf("%v:%v", ip, port)
	leolog.Print("Web Listener On:" + str)
	err := w.engine.Run(str)
	if err != nil {
		panic(fmt.Sprintf("Listener %s error: %s", str, err.Error()))
	}
}


func (engine *Engine) isProduction() bool {
	return engine.runMode == Pro
}

func (engine *Engine) InitLog(log Logger) {
	engine.log = log
}

func formatRequestLog(c *gin.Context) string {
	reqMethod := c.Request.Method
	reqUri := c.Request.RequestURI
	statusCode := c.Writer.Status()
	clientIP := c.Request.Host
	timeStr := leotime.DateFormat(leotime.TimeLocation(time.Now(), "Custom", 8), leotime.ForMate_yyyymmddhhmmss)
	str := fmt.Sprintf("| disTime: %s | code: %3d  | clientIP: %15s | reqMethod: %s | reqUri: %s |",
		timeStr, statusCode, clientIP, reqMethod, reqUri)
	return str
}
